package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"time"

	"github.com/a2659802/window-agent/pkg/agent"
	"github.com/a2659802/window-agent/pkg/message"
)

var (
	taskChan chan message.Message
)

var (
	queue = []message.Message{
		{
			Type: message.MsgTypePing,
		},
		{
			Type: message.MsgTypePong,
		},
		{
			Type: message.MsgTypePwdRequest,
			Data: agent.PasswordMessage{Action: agent.ActionChange, ID: "111", UserName: "test1", Password: "pwd1"},
		},
		{
			Type: message.MsgTypeResponse,
			Data: message.ResponseMessage{
				Code: message.StatusFailed,
				ID:   "111",
			},
		},
	}
)

func HandleClientConnect(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Receive Connect Request From ", conn.RemoteAddr().String())
	// buffer := make([]byte, 1024)
	for i := 0; i < len(queue); i += 2 {
		req := queue[i]
		expect := queue[i+1]

		log.Println("[info] get item from queue")
		if err := message.Encode(conn, &req); err != nil {
			log.Println(err.Error())
			break
		}
		log.Println("[info] send message to client")
		got, err := message.Decode(conn)
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("[info] got message from client:%v\n", string(got.RawData))
		if got.Type != expect.Type {
			log.Printf("unexpect, got %v, expect %v", got.Type, expect.Type)
			break
		}

	}
	time.Sleep(time.Second * 10)
	fmt.Println("Client " + conn.RemoteAddr().String() + " Connection Closed.....")
}

func main() {
	crtFile, keyFile := filepath.Join("ssl", "server.crt"), filepath.Join("ssl", "server.key")
	crt, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalln(err.Error())
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	// Time returns the current time as the number of seconds since the epoch.
	// If Time is nil, TLS uses time.Now.
	tlsConfig.Time = time.Now
	// Rand provides the source of entropy for nonces and RSA blinding.
	// If Rand is nil, TLS uses the cryptographic random reader in package
	// crypto/rand.
	// The Reader must be safe for use by multiple goroutines.
	tlsConfig.Rand = rand.Reader
	l, err := tls.Listen("tcp", ":8443", tlsConfig)
	if err != nil {
		log.Fatalln(err.Error())
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		} else {
			go HandleClientConnect(conn)
		}
	}

}
