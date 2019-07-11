package env

import (
	"fmt"
	"github.com/getgauge/common"
	"github.com/getgauge/gauge/env"
	"github.com/sitture/gauge-reportserver/logger"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

const (
	DefaultHost            = "http://localhost:8000"
	ReportServerHostEnv    = "REPORTSERVER_HOST"
	ReportServerBaseDirEnv = "REPORTSERVER_BASE_DIR"
	ReportServerPathEnv    = "REPORTSERVER_PATH"
	ReportServerTimeoutEnv = "REPORTSERVER_TIMEOUT_IN_SECONDS"
	LogsDirEnv             = "logs_directory"
	DefaultTimeout         = 15
)

func GetProjectRoot() string {
	return GetEnv(common.GaugeProjectRootEnv, true)
}

var GetProjectDirName = func() string {
	return path.Base(GetProjectRoot())
}

func GetCurrentExecutableDir() (string, string) {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf(err.Error())
	}
	target, err := filepath.EvalSymlinks(ex)
	if err != nil {
		return path.Dir(ex), filepath.Base(ex)
	}
	return filepath.Dir(target), filepath.Base(ex)
}

func GetReportsDir() (dir string) {
	dir = GetEnv(env.GaugeReportsDir, true)
	if filepath.IsAbs(dir) {
		return
	}
	dir = filepath.Join(GetProjectRoot(), dir)
	return
}

func GetEnv(envKey string, exitOnMissing bool) (value string) {
	value = os.Getenv(envKey)
	if len(value) == 0 && exitOnMissing {
		panic(fmt.Sprintf("Environment variable '%s' is not set. \n", envKey))
	}
	return
}

func GetEnvWithDefault(env, defaultValue string) (value string) {
	value = GetEnv(env, false)
	if len(value) == 0 {
		value = defaultValue
	}
	return
}

var PluginKillTimeout = func() int {
	value := GetEnv("plugin_kill_timeout", false)
	if len(value) == 0 {
		return 0
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return v / 1000
}

var ReportServerTimeout = func() time.Duration {
	timeout := GetEnvWithDefault(ReportServerTimeoutEnv, "15")
	int, err := strconv.ParseInt(timeout, 10, 64)
	if err != nil {
		logger.Infof("invalid value for '%s', setting to default", ReportServerTimeoutEnv)
		return time.Duration(DefaultTimeout) * time.Second
	}
	return time.Duration(int) * time.Second
}

var GaugeLogsFile = func() string {
	logsDir := GetEnvWithDefault(LogsDirEnv, "logs")
	return path.Join(GetProjectRoot(), logsDir, "gauge.log")
}

func GetReportServerHost() (url string) {
	url = GetEnv(ReportServerHostEnv, false)
	if len(url) == 0 {
		logger.Debugf("Could not find '%s', setting to default '%s'.", ReportServerHostEnv, DefaultHost)
		url = DefaultHost
	}
	return
}

func GetReportServerUrl() string {
	baseDir := GetEnvWithDefault(ReportServerBaseDirEnv, GetProjectDirName())
	logger.Debugf("baseDir => '%s'", baseDir)
	environment := GetEnvWithDefault(env.GaugeEnvironment, common.DefaultEnvDir)
	logger.Debugf("environment => '%s'", environment)
	reportPath := GetEnv(ReportServerPathEnv, false)
	logger.Debugf("reportPath => '%s'", reportPath)
	uri, err := url.Parse(GetReportServerHost())
	if err != nil {
		panic(err)
	}
	if len(reportPath) == 0 {
		uri.Path = path.Join(uri.Path, baseDir, environment, reportPath)
	} else {
		uri.Path = path.Join(uri.Path, baseDir, reportPath)
	}
	logger.Debugf("reportserverurl => '%s'", uri.String())
	return uri.String()
}
