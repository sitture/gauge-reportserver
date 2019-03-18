package env

import (
	"fmt"
	"github.com/getgauge/common"
	"github.com/haroon-sheikh/gauge/env"
	"log"
	"os"
	"path"
	"path/filepath"
)

func GetProjectRoot() string {
	return Getenv(common.GaugeProjectRootEnv)
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
	dir = Getenv(env.GaugeReportsDir)
	if filepath.IsAbs(dir) {
		return
	}
	dir = filepath.Join(GetProjectRoot(), dir)
	return
}

func Getenv(envKey string) (value string) {
	value = os.Getenv(envKey)
	if value == "" {
		fmt.Printf("Environment variable '%s' is not set. \n", envKey)
		os.Exit(1)
	}
	return
}

// TODO Look at env.go for parsing properties
