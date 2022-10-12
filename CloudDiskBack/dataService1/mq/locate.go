package mq

import (
	"CloudDIsk/dataService1/internal/config"
	"encoding/json"
)

// Locate 判断该节点是否存在
func Locate(node string) bool {
	if config.Conf.OBS.HuaweiObsBucketRootFolder == node {
		return true
	}
	return false
}

// StartLocate 监听交换机发送的消息
func StartLocate() {
	q := New(config.Conf.RabbitMQ.RabbitURL)
	defer q.Close()
	q.Bind(config.Conf.RabbitMQ.DataServers)
	c := q.Consume()

	for msg := range c {
		var o OBS
		err := json.Unmarshal(msg.Body, &o)
		if err != nil {
			panic(err)
		}
		if Locate(o.HuaweiObsBucketRootFolder) {
			q.Send(msg.ReplyTo, config.Conf.OBS)
		}
	}
}
