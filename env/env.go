package env

import (
	"fmt"
	"github.com/getgauge/common"
	"github.com/getgauge/gauge/env"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func GetProjectRoot() string {
	return GetEnv(common.GaugeProjectRootEnv)
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
	dir = GetEnv(env.GaugeReportsDir)
	if filepath.IsAbs(dir) {
		return
	}
	dir = filepath.Join(GetProjectRoot(), dir)
	return
}

func GetEnv(envKey string) (value string) {
	value = os.Getenv(envKey)
	if value == "" {
		fmt.Printf("Environment variable '%s' is not set. \n", envKey)
		os.Exit(1)
	}
	return
}

var PluginKillTimeout = func() int {
	value := GetEnv("plugin_kill_timeout")
	if value == "" {
		return 0
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return v / 1000
}

// TODO Look at env.go for parsing properties
