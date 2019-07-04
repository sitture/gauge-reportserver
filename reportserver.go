package main

import (
	"fmt"
	"github.com/sitture/gauge-reportserver/env"
	"github.com/sitture/gauge-reportserver/gauge_messages"
	"github.com/sitture/gauge-reportserver/listener"
	"github.com/sitture/gauge-reportserver/logger"
	"github.com/sitture/gauge-reportserver/sender"
	"github.com/sitture/gauge-reportserver/zipper"
	"os"
	"path"
	"strings"
	"time"
)

const (
	ReportServer      = "reportserver"
	PluginActionEnv   = ReportServer + "_action"
	ExecutionAction   = "execution"
	GaugeHost         = "127.0.0.1"
	GaugePortEnvVar   = "plugin_connection_port"
	HtmlReportDir     = "html-report"
	HtmlReportArchive = HtmlReportDir + ".zip"
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
		logger.Debug("Could not create the gauge listener")
		os.Exit(1)
	}
	shipper := &shipper{stopChan: stopChan}
	listener.OnSuiteResult(shipper.Send)
	listener.Start()
}

func (shipper *shipper) Send(suiteResult *gauge_messages.SuiteExecutionResult) {
	if IsReadyToShip() {
		SendReport(shipper.stopChan)
	}
}

func IsReadyToShip() (ready bool) {
	ready = false
	ticker := time.NewTicker(1 * time.Second)
	defer func() { ticker.Stop() }()
	timer := time.After(env.ReportServerTimeout())
	for {
		select {
		case <-timer:
			return
		case <-ticker.C:
			// do something
			if ReadLogsFile(env.GaugeLogsFile()) {
				fmt.Println("SENDING ...")
				return true
			}
		}
	}
	return
}

func ReadLogsFile(logsFilePath string) (logLineExists bool) {
	logLineExists = false
	logLine := "Plugin [Html Report] with pid"
	// check if logsFilePath exists
	if _, err := os.Stat(logsFilePath); os.IsNotExist(err) {
		logLineExists = false
	}
	file, err := os.Open(logsFilePath)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 512)
	stat, err := os.Stat(logsFilePath)
	start := stat.Size() - 512
	_, err = file.ReadAt(buf, start)
	if err == nil {
		logLineExists = strings.Contains(string(buf), logLine)
	}
	return
}

func SendReport(stop chan bool) {
	defer func(s chan bool) { s <- true }(stop)
	orig := path.Join(env.GetReportsDir(), HtmlReportDir)
	logger.Debug("Origin report directory is '%s'", orig)
	dest := path.Join(env.GetReportsDir(), HtmlReportArchive)
	logger.Debug("Archive destination is '%s'", dest)
	if err := zipper.ZipDir(orig, dest); err != nil {
		return
	}
	reportPath := env.GetReportServerUrl()
	err := sender.SendArchive(reportPath, dest)
	if err != nil {
		logger.Printf(fmt.Sprintf("Could not send the archive from '%s' to '%s'\n %s", dest, reportPath, err))
	}
	fmt.Printf("Successfully sent html-report to reportserver => %s", reportPath+"/report.html\n")
}
