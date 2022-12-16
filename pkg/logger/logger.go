package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/kardianos/service"
)

// 当以服务启动的时候，日志打到事件查看器中(eventvwr).
var logger service.Logger

func SetupLogger(s service.Service) {

	errs := make(chan error, 5)
	var err error
	logger, err = s.Logger(errs)
	if err != nil {
		log.Panicf("init logger error:%v", err)
		os.Exit(1)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warning(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	str := fmt.Sprintf("[Fatal] %s", fmt.Sprint(args...))
	Error(str)
	os.Exit(2)
}

func Fatalf(format string, args ...interface{}) {
	Errorf(format, args...)
	os.Exit(2)
}
