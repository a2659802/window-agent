// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// Simple service that only works by printing a log message every few seconds.
package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/a2659802/window-agent/pkg/agent"
	"github.com/a2659802/window-agent/pkg/config"
	"github.com/a2659802/window-agent/pkg/logger"
	"github.com/a2659802/window-agent/pkg/srvconn"
	"github.com/kardianos/service"
)

const (
	ServiceName        = "GoServiceExampleLogging"
	ServiceDisplayName = "Go Service Example for Logging"
	ServiceDescription = "This is an example Go service that outputs log messages."
)

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

// 处理服务启动
func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

// 真正的程序入口
func (p *program) run() {
	logger.Info("agent starting")

	// 加载配置文件
	config.Setup(*configFlag)

	// 获取服务地址
	addrStr := config.GlobalConfig.ServerHost

	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		logger.Fatal("cannot resolve server address")
	}
	// 与服务端建立连接
	tcpConn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		logger.Fatalf("connect fail:%s", err.Error())
	}

	// 创建context, 跟踪服务停止信号
	background := context.Background()
	ctx, cancel := context.WithCancel(background)
	defer cancel()

	go func() {
		// 等待服务停止信号
		<-p.exit
		cancel()
	}()

	cfg := config.GlobalConfig
	// 将连接所有权转交给srvconn模块
	conn := srvconn.NewConnection(tcpConn, cfg.ServerCAPath, cfg.ServerName)
	dispatchCh, toSendCh := conn.Start(ctx)

	// 启动agent
	a := agent.NewAgent(dispatchCh, toSendCh)
	if err := a.Run(ctx); err != nil {
		logger.Fatal(err.Error())
	}

}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}

// command-line option
var (
	svcFlag    = flag.String("service", "", "Control the system service.")
	configFlag = flag.String("config", "config.yaml", "config.yaml path")
)

// Service setup.
//   Define service config.
//   Create the service.
//   Setup the logger.
//   Handle service controls (optional).
//   Run the service.
func main() {
	flag.Parse()

	svcConfig := &service.Config{
		Name:        ServiceName,
		DisplayName: ServiceDisplayName,
		Description: ServiceDescription,
	}

	prg := &program{}

	// 注册服务
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// 初始化日志实例
	logger.SetupLogger(s)

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	// 阻塞直到Stop
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
