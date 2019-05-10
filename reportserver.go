package main

import (
	"bytes"
	"fmt"
	"github.com/haroon-sheikh/gauge-reportserver/env"
	"github.com/haroon-sheikh/gauge-reportserver/gauge_messages"
	"github.com/haroon-sheikh/gauge-reportserver/listener"
	"github.com/radovskyb/watcher"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	ReportServer    = "reportserver"
	PluginActionEnv = ReportServer + "_action"
	ExecutionAction = "execution"
	GaugeHost       = "127.0.0.1"
	GaugePortEnvVar = "plugin_connection_port"
)

var currentReportTimestamp = time.Now()

type shipper struct {
	result   *gauge_messages.SuiteExecutionResult
	stopChan chan bool
}

func ShipReport() {
	stopChan := make(chan bool)
	listener, err := listener.NewGaugeListener(GaugeHost, os.Getenv(GaugePortEnvVar), stopChan)
	if err != nil {
		fmt.Println("Could not create the gauge listener")
		os.Exit(1)
	}
	shipper := &shipper{stopChan: stopChan}
	listener.OnSuiteStart(printme2)
	listener.OnSuiteEnd(printme3)
	listener.OnSuiteResult(shipper.Meta)
	listener.OnKill(printme4)
	listener.Start()
}

func (shipper *shipper) Meta(suiteResult *gauge_messages.SuiteExecutionResult) {
	SendReport(shipper.stopChan)
}

func SendReport(stop chan bool) {
	defer func(s chan bool) { s <- true }(stop)
	fmt.Println("HELLOOOOOOOOO1")
	fmt.Printf("Successfully sent html-report to reportserver => ")
}

func sendme(stop chan bool) {
	defer func(s chan bool) { s <- true }(stop)
	SendReport()
}

func printme2(suiteResult *gauge_messages.ExecutionStartingRequest) {
	fmt.Println("HELLO, ExecutionStartingRequest!")
	fmt.Println(getMod().After(currentReportTimestamp))
}

func printme4(suiteResult *gauge_messages.KillProcessRequest) {
	fmt.Println("HELLO, KillProcessRequest!")
	//getMod()
}

func printme3(suiteResult *gauge_messages.ExecutionEndingRequest) {
	fmt.Println("HELLO, ExecutionEndingRequest!")
	//getMod()
}

func (shipper *shipper) printme(suiteResult *gauge_messages.SuiteExecutionResult) {
	//fmt.Println(suiteResult.GetSuiteResult().GetEnvironment())
	//fmt.Println(suiteResult.GetSuiteResult())
	//fmt.Println("HELLO, SuiteExecutionResult!")
	//dir, _ := os.Getwd()
	//fmt.Println(dir)
	//time.Sleep(30 * time.Second)
	//fmt.Println(getMod().After(currentReportTimestamp))
	////stdout()
	////fmt.Println(std2())
	//capture()
	//a := env.GetReportsDir()
	//fmt.Println("Reports Dir: " + a)
	//fmt.Println(IsReportGenerated())
	//fmt.Println(getMod().After(currentReportTimestamp))
	//
	sendme(r.stopChan)
	fmt.Println("HELLOOOOOOOOO1")
	//var maxNoUpdates = 0
	//for {
	//	if maxNoUpdates < 5 {
	//		if getMod().After(currentReportTimestamp) {
	//			time.Sleep(100 * time.Millisecond)
	//			maxNoUpdates = 0
	//		}
	//		maxNoUpdates++
	//		continue
	//	}
	//	return
	//}
}

func IsReportGenerated() (generated bool) {

	w := watcher.New()
	defer w.Close()
	w.FilterOps(watcher.Write)
	w.SetMaxEvents(1)

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event.Path) // Print the event's info.
				done <- true
				return
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	go func() {
		defer func() { done <- true }()
		if <-done {
			w.Wait()
			w.Close()
		}
	}()

	// Watch this folder for changes.
	if err := w.AddRecursive("/Users/has23/workspace/id/europa-e2e/reports/"); err != nil {
		log.Fatalln(err)
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
	<-done

	return true
}

func getMod() time.Time {
	if _, err := os.Stat(env.GetProjectRoot() + "/reports"); !os.IsNotExist(err) {
		// path/to/whatever exists
		file, _ := os.Stat(env.GetProjectRoot() + "/reports/html-report/index.html")
		fmt.Println(file.ModTime())
		return file.ModTime()
	}
	return currentReportTimestamp
}

func std2() string {

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()

}

func stdout() {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	print()
	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	// reading our temp stdout
	fmt.Println("---------")
	fmt.Println(out)
}

func capture() {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	fmt.Printf("Captured: %s", out) // prints: Captured: Hello, playground
}

//func main() {
//	fmt.Println(os.Getwd())
//	timea, _ := os.Stat("main.go")
//	mod := timea.ModTime()
//	a, _:= time.Parse(time.RFC822, mod.String())
//	fmt.Println(a)
//
//	fmt.Println(time.Now().After(mod))
//
//}
