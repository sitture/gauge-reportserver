package env

import (
	"github.com/getgauge/common"
	"github.com/getgauge/gauge/env"
	"os"
	"testing"
	"time"
)

func TestGetProjectRoot(t *testing.T) {
	expected := "/test/test-project"
	os.Setenv(common.GaugeProjectRootEnv, expected)
	defer func() { os.Unsetenv(common.GaugeProjectRootEnv) }()
	projectRoot := GetProjectRoot()
	if projectRoot != expected {
		t.Errorf("GetProjectRoot was incorrect, got: %s, want: %s.", projectRoot, expected)
	}
}

func TestGetProjectDirName(t *testing.T) {
	expected := "test-project"
	os.Setenv(common.GaugeProjectRootEnv, "/test/"+expected)
	defer func() { os.Unsetenv(common.GaugeProjectRootEnv) }()
	projectDir := GetProjectDirName()
	if projectDir != expected {
		t.Errorf("GetProjectDirName was incorrect, got: %s, want: %s.", projectDir, expected)
	}
}

func TestGetReportsDir(t *testing.T) {
	os.Setenv(common.GaugeProjectRootEnv, "/test/test-project")
	os.Setenv(env.GaugeReportsDir, "reports")
	defer func() {
		os.Unsetenv(common.GaugeProjectRootEnv)
		os.Unsetenv(env.GaugeReportsDir)
	}()
	reportsDir := GetReportsDir()
	expected := "/test/test-project/reports"
	if reportsDir != expected {
		t.Errorf("GetReportsDir was incorrect, got: %s, want: %s.", reportsDir, expected)
	}
}

func TestGaugeLogsFile(t *testing.T) {
	os.Setenv(common.GaugeProjectRootEnv, "/test/test-project")
	defer func() { os.Unsetenv(common.GaugeProjectRootEnv) }()

	logsFile := GaugeLogsFile()
	expected := "/test/test-project/logs/gauge.log"
	if logsFile != expected {
		t.Errorf("GaugeLogsFile was incorrect, got: %s, want: %s.", logsFile, expected)
	}
}

func TestReportServerTimeout(t *testing.T) {
	timeout := ReportServerTimeout()
	expected := DefaultTimeout * time.Second
	if timeout != expected {
		t.Errorf("ReportServerTimeout was incorrect, got: %s, want: %s.", timeout, expected)
	}
}

func TestReportServerTimeoutCustom(t *testing.T) {
	os.Setenv(ReportServerTimeoutEnv, "20")
	defer func() { os.Unsetenv(ReportServerTimeoutEnv) }()
	timeout := ReportServerTimeout()
	expected := 20 * time.Second
	if timeout != expected {
		t.Errorf("ReportServerTimeout was incorrect, got: %s, want: %s.", timeout, expected)
	}
}

func TestReportServerTimeoutInvalid(t *testing.T) {
	os.Setenv(ReportServerTimeoutEnv, "invalid")
	defer func() { os.Unsetenv(ReportServerTimeoutEnv) }()
	timeout := ReportServerTimeout()
	expected := DefaultTimeout * time.Second
	if timeout != expected {
		t.Errorf("ReportServerTimeout was incorrect, got: %s, want: %s.", timeout, expected)
	}
}

func TestGetReportServerHostDefault(t *testing.T) {
	reportServerUrl := GetReportServerHost()
	if reportServerUrl != DefaultHost {
		t.Errorf("GetReportServerHost was incorrect, got: %s, want: %s.", reportServerUrl, DefaultHost)
	}
}

func TestGetReportServerHost(t *testing.T) {
	expected := "http://testing:8080"
	os.Setenv(ReportServerHostEnv, expected)
	defer func() { os.Unsetenv(ReportServerHostEnv) }()
	reportServerUrl := GetReportServerHost()
	if reportServerUrl != expected {
		t.Errorf("GetReportServerHost was incorrect, got: %s, want: %s.", reportServerUrl, expected)
	}
}

func TestGetReportServerUrlDefaultEnv(t *testing.T) {
	os.Setenv(common.GaugeProjectRootEnv, "/test/test-project")
	defer func() { os.Unsetenv(common.GaugeProjectRootEnv) }()
	reportServerUrl := GetReportServerUrl()
	expected := "http://localhost:8000/test-project/default"
	if reportServerUrl != expected {
		t.Errorf("GetReportServerUrl was incorrect, got: %s, want: %s.", reportServerUrl, expected)
	}
}

func TestGetReportServerUrlCustomEnv(t *testing.T) {
	os.Setenv(common.GaugeProjectRootEnv, "/test/test-project")
	os.Setenv(env.GaugeEnvironment, "test")
	defer func() {
		os.Unsetenv(common.GaugeProjectRootEnv)
		os.Unsetenv(env.GaugeEnvironment)
	}()
	reportServerUrl := GetReportServerUrl()
	expected := "http://localhost:8000/test-project/test"
	if reportServerUrl != expected {
		t.Errorf("GetReportServerUrl was incorrect, got: %s, want: %s.", reportServerUrl, expected)
	}
}

func TestGetReportServerUrlBaseDir(t *testing.T) {
	os.Setenv(common.GaugeProjectRootEnv, "/test/test-project")
	os.Setenv(ReportServerBaseDirEnv, "test")
	defer func() {
		os.Unsetenv(common.GaugeProjectRootEnv)
		os.Unsetenv(ReportServerBaseDirEnv)
	}()
	reportServerUrl := GetReportServerUrl()
	expected := "http://localhost:8000/test/default"
	if reportServerUrl != expected {
		t.Errorf("GetReportServerUrl was incorrect, got: %s, want: %s.", reportServerUrl, expected)
	}
}

func TestGetReportServerUrlBaseDirAndPath(t *testing.T) {
	os.Setenv(common.GaugeProjectRootEnv, "/test/test-project")
	os.Setenv(ReportServerBaseDirEnv, "test")
	os.Setenv(ReportServerPathEnv, "hello/world")
	defer func() {
		os.Unsetenv(common.GaugeProjectRootEnv)
		os.Unsetenv(ReportServerBaseDirEnv)
		os.Unsetenv(ReportServerPathEnv)
	}()
	reportServerUrl := GetReportServerUrl()
	expected := "http://localhost:8000/test/hello/world"
	if reportServerUrl != expected {
		t.Errorf("GetReportServerUrl was incorrect, got: %s, want: %s.", reportServerUrl, expected)
	}
}
