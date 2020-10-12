package listener

import (
	"bytes"
	"fmt"
	"github.com/getgauge/common"
	"github.com/golang/protobuf/proto"
	"github.com/sitture/gauge-reportserver/env"
	"github.com/sitture/gauge-reportserver/gauge_messages"
	"github.com/sitture/gauge-reportserver/logger"
	"net"
	"os"
	"time"
)

type GaugeResultHandlerFn func(result *gauge_messages.SuiteExecutionResult)
type KillProcessRequestHandlerFn func(killProcessRequest *gauge_messages.KillProcessRequest)

type Listener struct {
	connection                  net.Conn
	onResultHandler             GaugeResultHandlerFn
	onKillProcessRequestHandler KillProcessRequestHandlerFn
	stopChan                    chan bool
}

func NewGaugeListener(host string, port string, stopChan chan bool) (*Listener, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err == nil {
		return &Listener{connection: conn, stopChan: stopChan}, nil
	}
	return nil, err
}

func (listener *Listener) OnSuiteResult(resultHandler GaugeResultHandlerFn) {
	listener.onResultHandler = resultHandler
}

func (listener *Listener) OnKillProcessRequest(killProcessRequestHandler KillProcessRequestHandlerFn) {
	listener.onKillProcessRequestHandler = killProcessRequestHandler
}

func (listener *Listener) Start() {
	buffer := new(bytes.Buffer)
	data := make([]byte, 8192)
	for {
		n, err := listener.connection.Read(data)
		if err != nil {
			return
		}
		buffer.Write(data[0:n])
		listener.ProcessMessages(buffer)
	}
}

func (listener *Listener) ProcessMessages(buffer *bytes.Buffer) {
	for {
		messageLength, bytesRead := proto.DecodeVarint(buffer.Bytes())
		if messageLength <= 0 || messageLength >= uint64(buffer.Len()) {
			return
		}
		message := &gauge_messages.Message{}
		messageBoundary := int(messageLength) + bytesRead
		messageBytes := buffer.Bytes()[bytesRead:messageBoundary]
		if err := proto.Unmarshal(messageBytes, message); err != nil {
			logger.Warnf("Failed to read proto message: %s\n", err.Error())
			logger.Warnf("Message : %s\n", string(messageBytes))
			continue
		}
		switch message.MessageType {
		case gauge_messages.Message_KillProcessRequest:
			logger.Debug("Received Kill Message, exiting...")
			listener.onKillProcessRequestHandler(message.GetKillProcessRequest())
			err := listener.connection.Close()
			if err != nil {
				logger.Debug("Failed to close the listener connection.")
			}
			os.Exit(0)
		case gauge_messages.Message_SuiteExecutionResult:
			go listener.sendPings()
			listener.onResultHandler(message.GetSuiteExecutionResult())
		}
		buffer.Next(messageBoundary)
		if buffer.Len() == 0 {
			return
		}
	}
}

func (listener *Listener) sendPings() {
	msg := &gauge_messages.Message{
		MessageId:   common.GetUniqueID(),
		MessageType: gauge_messages.Message_KeepAlive,
		KeepAlive:   &gauge_messages.KeepAlive{PluginId: "reportserver"},
	}
	m, err := proto.Marshal(msg)
	if err != nil {
		logger.Debugf("Unable to marshal ping message, %s", err.Error())
		return
	}
	ping := func(b []byte, c net.Conn) {
		logger.Debug("reportserver sending a keep-alive ping")
		l := proto.EncodeVarint(uint64(len(b)))
		if _, err := c.Write(append(l, b...)); err != nil {
			logger.Debugf("Unable to send ping message, %s", err.Error())
		}
	}
	ticker := time.NewTicker(interval())
	defer ticker.Stop()

	for {
		select {
		case <-listener.stopChan:
			logger.Debug("Stopping pings")
			return
		case <-ticker.C:
			ping(m, listener.connection)
		}
	}
}

var interval = func() time.Duration {
	v := env.PluginKillTimeout()
	if v/2 < 2 {
		return 2 * time.Second
	}
	return time.Duration(v * 1000 * 1000 * 1000 / 2)
}
