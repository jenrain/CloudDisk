package mq

import (
	"core/internal/config"
	"strconv"
	"time"
)

// Locate 定位文件
func Locate(node string) string {
	q := New(config.Conf.RabbitMQ.RabbitURL, "")
	q.Publish(config.Conf.RabbitMQ.DataServers, node)
	c := q.Consume()

	// 一秒后关闭临时消息队列
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))

	return s
}

func Exist(name string) bool {
	return Locate(name) != ""
}
