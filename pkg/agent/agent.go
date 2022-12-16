package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/a2659802/window-agent/pkg/logger"
	"github.com/a2659802/window-agent/pkg/message"
)

// Agent 不感知具体的网络通信, 使用go channel来与上一层模块交换数据
type Agent struct {
	handlers map[message.MessageType]func(data []byte)
	recvCh   <-chan message.Message // 堡垒机发过来的消息通道
	sendCh   chan<- message.Message // agent 需要发给堡垒机的

	handleStatistic int
}

func NewAgent(recv <-chan message.Message, send chan<- message.Message) *Agent {
	if recv == nil || send == nil {
		logger.Fatal("init agent fail: cannot use nil channel to initialize")
	}
	return &Agent{
		handlers: make(map[message.MessageType]func(data []byte)),
		recvCh:   recv,
		sendCh:   send,
	}
}

func (a *Agent) Run(ctx context.Context) error {
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case msg, ok := <-a.recvCh:
			if !ok {
				return fmt.Errorf("error read message from channel: closed")
			}
			a.handle(msg)
		}
	}

	return nil
}

// 处理消息
func (a *Agent) handle(msg message.Message) {
	handlers := map[message.MessageType]func(data []byte){
		// PING -> PONG
		message.MsgTypePing: func(data []byte) { a.pong() },
		// CHANGE PASSWORD -> RESPONSE
		// VERIFY PASSWORD -> RESPONSE
		message.MsgTypePwdRequest: a.processPwdMessage,
	}
	h, ok := handlers[msg.Type]
	if !ok {
		a.response("", message.StatusBadRequest, "unknown message type")
		return
	}

	h(msg.RawData)
	a.handleStatistic++
}

func (a *Agent) response(id string, code message.StatusCode, reasons ...string) {
	reason := strings.Join(reasons, ";")

	msg := message.Message{
		Type: message.MsgTypeResponse,
		Data: message.ResponseMessage{
			ID:     id,
			Code:   code,
			Reason: reason,
		},
	}

	a.sendCh <- msg
}

func (a *Agent) pong() {
	msg := message.Message{
		Type: message.MsgTypePong,
	}

	a.sendCh <- msg
}

func (a *Agent) AddRoute(key message.MessageType, handler func(data []byte)) {
	a.handlers[key] = handler
}
