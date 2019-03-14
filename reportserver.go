package main

import (
	"fmt"
	"github.com/haroon-sheikh/gauge-report-server/gauge_messages"
	"github.com/haroon-sheikh/gauge-report-server/listener"
	"os"
)

const (
	ReportServer    = "report-server"
	PluginActionEnv = "report-server_action"
	ExecutionAction = "execution"
	GaugeHost       = "127.0.0.1"
	GaugePortEnvVar = "plugin_connection_port"
)

func sendReport() {
	listener, err := listener.NewGaugeListener(GaugeHost, os.Getenv(GaugePortEnvVar))
	if err != nil {
		fmt.Println("Could not create the gauge listener")
		os.Exit(1)
	}
	listener.OnSuiteStart(printme2)
	listener.OnSuiteResult(printme)
	listener.Start()
}

func printme2(suiteResult *gauge_messages.ExecutionStartingRequest) {
	fmt.Println("HELLO, ExecutionStartingRequest!")
}

func printme(suiteResult *gauge_messages.SuiteExecutionResult) {
	fmt.Println(suiteResult.GetSuiteResult().GetEnvironment())
	fmt.Println(suiteResult.GetSuiteResult())
	fmt.Println("HELLO, SuiteExecutionResult!")
}
