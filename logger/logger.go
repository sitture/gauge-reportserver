package logger

import (
	"log"
	"os"
	"strings"
)

var debug bool

func Init() {
	var logLevel = os.Getenv("GAUGE_LOG_LEVEL")
	if strings.ToLower(logLevel) == "debug" {
		debug = true
	}
	log.SetPrefix("[reportserver] ")
}

func Debug(format string, args ...string) {
	if debug {
		log.Printf(format, args)
	}
}

func Printf(format string, args ...string) {
	log.Printf(format, args)
}
