package config

import (
	"github.com/a2659802/window-agent/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	ServerHost   string `validate:"required,hostname_port"` //  服务端地址
	ServerCAPath string `validate:"required,file"`          //  服务端证书路径
	ServerName   string `validate:"required,hostname"`      //  服务名
	LogLevel     string
}

var GlobalConfig *Config

func (c *Config) EnsureConfigValid() {
	// 根据tag自动校验配置
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		logger.Fatalf("verify config fail:%v", err.Error())
	}
}

func Setup(configPath string) {
	var conf = getDefaultConfig()
	loadConfigFromFile(configPath, &conf)
	conf.EnsureConfigValid()
	GlobalConfig = &conf
}

func getDefaultConfig() Config {
	return Config{}
}

func loadConfigFromFile(path string, conf *Config) {
	fileViper := viper.New()
	fileViper.SetConfigFile(path)
	if err := fileViper.ReadInConfig(); err != nil {
		logger.Errorf("fail to load config file:%v", path)
	}
	if err := fileViper.Unmarshal(conf); err != nil {
		logger.Errorf("fail to parse config file:%v", path)
		return
	}

	logger.Infof("load config file %v success", path)
}
