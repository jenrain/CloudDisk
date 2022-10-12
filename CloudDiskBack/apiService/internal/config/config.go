package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Password string
		Addr     string
	}
	RabbitMQ struct {
		AsyncTransferEnable bool
		RabbitURL           string
		ApiServers          string
		DataServers         string
		CanalExchange       string
	}
}

var Conf Config
