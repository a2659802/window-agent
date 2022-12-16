package srvconn

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"sync"

	"github.com/a2659802/window-agent/pkg/logger"
	"github.com/a2659802/window-agent/pkg/message"
)

const (
	MsgChanSize = 10
)

type Connection struct {
	sync.Mutex
	conn       net.Conn
	tlsConn    *tls.Conn
	ca         []byte
	serverName string

	sendStatistic int // 发送消息计数器
	recvStatistic int // 接收消息计数器
}

func NewConnection(conn net.Conn, caPath, sni string) *Connection {
	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		logger.Fatalf("load ca file error:%v", err.Error())
	}

	return &Connection{
		conn:       conn,
		ca:         caCert,
		serverName: sni,
	}
}

func (c *Connection) Start(ctx context.Context) (<-chan message.Message, chan<- message.Message) {
	// 加载根证书
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(c.ca)
	config := &tls.Config{
		// Certificates:       []tls.Certificate{cliCert}, // 客户端证书, 双向认证必须携带
		RootCAs:            caCertPool, // 校验服务端证书 [CA证书]
		InsecureSkipVerify: false,      // 不用校验服务器证书
		ServerName:         c.serverName,
	}

	// tls over tcp
	tlsConn := tls.Client(c.conn, config)

	// 握手
	if err := tlsConn.HandshakeContext(ctx); err != nil {
		logger.Errorf("handshake error:%v", err.Error())
		c.conn.Close()
		return nil, nil
	}

	c.tlsConn = tlsConn

	go func() {
		<-ctx.Done()
		// break pipe
		c.Close()
	}()

	dispatchCh, toSendCh := make(chan message.Message, MsgChanSize), make(chan message.Message, MsgChanSize)

	go c.readLoop(dispatchCh)
	go c.writeLoop(toSendCh)

	return dispatchCh, toSendCh
}

// 循环读取服务端数据，将数据切割成消息结构，发往消息处理通道
func (c *Connection) readLoop(msgCh chan<- message.Message) {
	// 主动退出程序或者连接断开，关闭消息通道，让上层应用感知
	defer close(msgCh)

	for {
		// 消息解码
		msg, err := message.Decode(c)
		if err != nil {
			logger.Error(err.Error())
			break
		}
		c.recvStatistic++

		// 发送到应用层通道，交给应用处理
		msgCh <- *msg
	}

	logger.Info("connect end of reading")
}

func (c *Connection) writeLoop(msgCh <-chan message.Message) {
	for {
		// 从待发送通道取消息
		msg, ok := <-msgCh
		if !ok {
			logger.Info("write loop break: channel close")
			break
		}

		// 消息编码
		buf := &bytes.Buffer{}
		if err := message.Encode(buf, &msg); err != nil {
			logger.Errorf("fail to encode message:%v", err.Error())
			break
		}
		// 发送消息
		nw, err := c.Write(buf.Bytes())
		if err != nil {
			logger.Error("fail to write to peer")
			break
		}

		if nw != buf.Len() {
			logger.Errorf("write not complete, write %v bytes", nw)
			break
		}

		c.sendStatistic++
	}

	logger.Info("connect end of writing")
}

func (c *Connection) Read(buf []byte) (int, error) {
	return c.tlsConn.Read(buf)
}

func (c *Connection) Write(buf []byte) (int, error) {
	return c.tlsConn.Write(buf)
}

func (c *Connection) Close() error {
	c.Lock()
	defer c.Unlock()

	if c.tlsConn != nil {
		if err := c.tlsConn.Close(); err != nil {
			return err
		}
		c.tlsConn = nil
	}

	return nil
}
