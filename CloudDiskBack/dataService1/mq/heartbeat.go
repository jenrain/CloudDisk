package mq

import (
	"CloudDIsk/dataService1/internal/config"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"time"
)

// StartHeartbeat 每隔 5s 发送心跳
func StartHeartbeat() {
	q := New(config.Conf.RabbitMQ.RabbitURL)
	defer q.Close()

	for {
		if isNodeHealthy() {
			q.Publish(config.Conf.RabbitMQ.ApiServers, config.Conf.OBS)
		}
		time.Sleep(5 * time.Second)
	}
}

func isNodeHealthy() bool {
	var obsClient, err = obs.New(config.Conf.OBS.HuaweiObsAK, config.Conf.OBS.HuaweiObsSK, config.Conf.OBS.HuaweiObsEndPoint)
	if err != nil {
		fmt.Println("obsClient 结构体创建失败")
	}
	_, err = obsClient.HeadBucket(config.Conf.OBS.HuaweiObsBucket)
	if err == nil {
		return true
	}
	return false
}
