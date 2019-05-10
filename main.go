package main

import "os"

func main() {
	action := os.Getenv(PluginActionEnv)
	if action == ExecutionAction {
		ShipReport()
	}
}
