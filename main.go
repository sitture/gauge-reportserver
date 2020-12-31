package main

import (
	"github.com/sitture/gauge-reportserver/gauge_messages"
	"github.com/sitture/gauge-reportserver/logger"
	"google.golang.org/grpc"
	"net"
	"os"
)

const oneGB = 1024 * 1024 * 1024

func main() {
	if os.Getenv(PluginActionEnv) == ExecutionAction {
		address, err := net.ResolveTCPAddr("tcp", GaugeHost)
		if err != nil {
			logger.Fatal("failed to start server.")
		}
		listener, err := net.ListenTCP("tcp", address)
		if err != nil {
			logger.Fatal("failed to start server.")
		}
		server := grpc.NewServer(grpc.MaxRecvMsgSize(oneGB))
		h := &handler{server: server}
		gauge_messages.RegisterReporterServer(server, h)
		logger.Infof("Listening on port:%d", listener.Addr().(*net.TCPAddr).Port)
		server.Serve(listener)
	}
}
