package listener

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/haroon-sheikh/gauge-reportserver/gauge_messages"
	"github.com/haroon-sheikh/gauge-reportserver/logger"
)

type GaugeSuiteStartHandlerFn func(result *gauge_messages.ExecutionStartingRequest)
type GaugeSuiteEndHandlerFn func(result *gauge_messages.ExecutionEndingRequest)
type GaugeResultHandlerFn func(result *gauge_messages.SuiteExecutionResult)
type GaugeKillProcessHandlerFn func(result *gauge_messages.KillProcessRequest)

type Listener struct {
	connection          net.Conn
	onResultHandler     GaugeResultHandlerFn
	onSuiteStartHandler GaugeSuiteStartHandlerFn
	onSuiteEndHandler   GaugeSuiteEndHandlerFn
	onKillHander        GaugeKillProcessHandlerFn
}

func NewGaugeListener(host string, port string) (*Listener, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err == nil {
		return &Listener{connection: conn}, nil
	}
	return nil, err
}

func (listener *Listener) OnSuiteStart(resultHandler GaugeSuiteStartHandlerFn) {
	listener.onSuiteStartHandler = resultHandler
}

func (listener *Listener) OnSuiteEnd(resultHandler GaugeSuiteEndHandlerFn) {
	listener.onSuiteEndHandler = resultHandler
}

func (listener *Listener) OnSuiteResult(resultHandler GaugeResultHandlerFn) {
	listener.onResultHandler = resultHandler
}

func (listener *Listener) OnKill(resultHandler GaugeKillProcessHandlerFn) {
	listener.onKillHander = resultHandler
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
		if messageLength > 0 && messageLength < uint64(buffer.Len()) {
			message := &gauge_messages.Message{}
			messageBoundary := int(messageLength) + bytesRead
			err := proto.Unmarshal(buffer.Bytes()[bytesRead:messageBoundary], message)
			if err != nil {
				log.Printf("Failed to read proto message: %s\n", err.Error())
			} else {
				switch message.MessageType {
				case gauge_messages.Message_KillProcessRequest:
					logger.Debug("Received Kill Message, exiting...")
					listener.onKillHander(message.GetKillProcessRequest())
					time.Sleep(9 * time.Second)
					listener.connection.Close()
					os.Exit(0)
				case gauge_messages.Message_ExecutionStarting:
					listener.onSuiteStartHandler(message.GetExecutionStartingRequest())
				case gauge_messages.Message_ExecutionEnding:
					listener.onSuiteEndHandler(message.GetExecutionEndingRequest())
				case gauge_messages.Message_SuiteExecutionResult:
					listener.onResultHandler(message.GetSuiteExecutionResult())
				}
				buffer.Next(messageBoundary)
				if buffer.Len() == 0 {
					return
				}
			}
		} else {
			return
		}
	}
}
