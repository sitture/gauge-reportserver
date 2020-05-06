package main

import (
	"github.com/sitture/gauge-reportserver/env"
	"github.com/sitture/gauge-reportserver/gauge_messages"
	"github.com/sitture/gauge-reportserver/listener"
	"github.com/sitture/gauge-reportserver/logger"
	"github.com/sitture/gauge-reportserver/sender"
	"github.com/sitture/gauge-reportserver/zipper"
	"os"
	"path"
	"path/filepath"
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
	// OldIndexFilePath is the name of index file
	OldIndexFile = "index.html"
	// NewIndexFilePath is the name of the new index file
	NewIndexFile = "report.html"
)

type shipper struct {
	result   *gauge_messages.SuiteExecutionResult
	stopChan chan bool
}

func ShipReport() {
	stopChan := make(chan bool)
	gaugeListener, err := listener.NewGaugeListener(GaugeHost, os.Getenv(GaugePortEnvVar), stopChan)
	if err != nil {
		logger.Fatal("Could not create the gauge listener")
		os.Exit(1)
	}
	shipper := &shipper{stopChan: stopChan}
	gaugeListener.OnSuiteResult(shipper.Send)
	gaugeListener.OnKillProcessRequest(shipper.TearDown)
	gaugeListener.Start()
}

func (shipper *shipper) Send(suiteResult *gauge_messages.SuiteExecutionResult) {
	if IsReadyToShip() {
		SendReport(shipper.stopChan)
	}
}

func (shipper *shipper) TearDown(killRequest *gauge_messages.KillProcessRequest) {
	// Check and delete existing archive
	if err := RemoveExistingArchive(ArchiveDestination()); err != nil {
		logger.Infof("Could not remove archive '%s'.", ArchiveDestination())
	}
	// Rename report.html back to index.html
	if err := RenameIndexFile(ArchiveOrigin(), NewIndexFile, OldIndexFile); err != nil {
		logger.Infof("Could not rename file from '%s' to '%s'.", NewIndexFile, OldIndexFile)
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
			logger.Infof("html-report was not ready, Timed out!")
			return
		case <-ticker.C:
			return ReadLogsFile(env.GaugeLogsFile())
		}
	}
}

func ReadLogsFile(logsFilePath string) (logLineExists bool) {
	logLineExists = false
	logLine := "Done generating HTML report"
	// check if logsFilePath exists
	if _, err := os.Stat(logsFilePath); os.IsNotExist(err) {
		logLineExists = false
	}
	file, err := os.Open(logsFilePath)
	defer func() {
		err := file.Close()
		if err != nil {
			return
		}
	}()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 512)
	stat, err := os.Stat(logsFilePath)
	if err != nil {
		return
	}
	start := stat.Size() - 512
	_, err = file.ReadAt(buf, start)
	if err == nil {
		logLineExists = strings.Contains(string(buf), logLine)
	}
	return
}

var ArchiveDestination = func() (dest string) {
	dest = path.Join(env.GetReportsDir(), HtmlReportArchive)
	logger.Debugf("Archive destination is '%s'", dest)
	return
}

var ArchiveOrigin = func() (orig string) {
	orig = path.Join(env.GetReportsDir(), HtmlReportDir)
	logger.Debugf("Origin report directory is '%s'", orig)
	return
}

func SendReport(stop chan bool) {
	defer func(s chan bool) { s <- true }(stop)
	// Rename index.html to report.html
	if err := RenameIndexFile(ArchiveOrigin(), OldIndexFile, NewIndexFile); err != nil {
		logger.Infof("Could not rename file from '%s' to '%s'\n", OldIndexFile, NewIndexFile)
	}
	// Check and delete existing archive
	if err := RemoveExistingArchive(ArchiveDestination()); err != nil {
		logger.Infof("Could not remove archive '%s'\n", ArchiveDestination())
	}
	if err := zipper.ZipDir(ArchiveOrigin(), ArchiveDestination()); err != nil {
		logger.Infof("error archiving the reports directory.\n%s\n", err.Error())
		return
	}
	reportPath := env.GetReportServerUrl()
	err := sender.SendArchive(reportPath, ArchiveDestination())
	if err != nil {
		logger.Infof("Could not send the archive from '%s' to '%s'\n%s\n", ArchiveDestination(), reportPath, err.Error())
	} else {
		logger.Infof("Successfully sent html-report to reportserver => %s\n", filepath.Join(reportPath, "report.html"))
	}
}

func RenameIndexFile(dir, from, to string) (err error) {
	logger.Debugf("renaming index file to '%s'", to)
	if _, err := os.Stat(path.Join(dir, from)); err == nil {
		err = os.Rename(path.Join(dir, from), path.Join(dir, to))
	}
	return
}

func RemoveExistingArchive(archivePath string) (err error) {
	logger.Debugf("removing archive '%s'", archivePath)
	if _, err := os.Stat(archivePath); err == nil {
		err = os.Remove(archivePath)
	}
	return
}
