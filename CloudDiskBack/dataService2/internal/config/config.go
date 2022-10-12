package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	RabbitMQ struct {
		AsyncTransferEnable bool
		RabbitURL           string
		ApiServers          string
		DataServers         string
	}
	OBS struct {
		HuaweiObsEndPoint         string
		HuaweiObsDawnLoadUrl      string
		HuaweiObsBucket           string
		HuaweiObsBucketRootFolder string
		HuaweiObsAK               string
		HuaweiObsSK               string
	}
}

var Conf Config
