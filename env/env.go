package env

import (
	"fmt"
	"github.com/getgauge/common"
	"github.com/getgauge/gauge/env"
	"github.com/haroon-sheikh/gauge-reportserver/logger"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

const (
	DefaultHost            = "http://localhost:8000"
	ReportServerHostEnv    = "REPORTSERVER_HOST"
	ReportServerBaseDirEnv = "REPORTSERVER_BASE_DIR"
	ReportServerPathEnv    = "REPORTSERVER_PATH"
	// GaugeEnvironmentEnv holds the name of the current environment
	GaugeEnvironmentEnv = "gauge_environment"
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
	if value == "" && exitOnMissing {
		panic(fmt.Sprintf("Environment variable '%s' is not set. \n", envKey))
	}
	return
}

func GetEnvWithDefault(env, defaultValue string) (value string) {
	value = GetEnv(env, false)
	if value == "" {
		value = defaultValue
	}
	return
}

var PluginKillTimeout = func() int {
	value := GetEnv("plugin_kill_timeout", false)
	if value == "" {
		return 0
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return v / 1000
}

func GetReportServerHost() (url string) {
	url = GetEnv(ReportServerHostEnv, false)
	if url == "" {
		logger.Debug("Could not find '%s', setting to default '%s'.", ReportServerHostEnv, DefaultHost)
		url = DefaultHost
	}
	return
}

func GetReportServerUrl() string {
	baseDir := GetEnvWithDefault(ReportServerBaseDirEnv, GetProjectDirName())
	environment := GetEnvWithDefault(GaugeEnvironmentEnv, common.DefaultEnvDir)
	reportPath := GetEnv(ReportServerPathEnv, false)
	uri, err := url.Parse(GetReportServerHost())
	if err != nil {
		panic(err)
	}
	if reportPath == "" {
		uri.Path = path.Join(uri.Path, baseDir, environment, reportPath)
	} else {
		uri.Path = path.Join(uri.Path, baseDir, reportPath)
	}
	return uri.String()
}
