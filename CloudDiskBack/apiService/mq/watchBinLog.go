package mq

import (
	"core/internal/config"
	"core/models"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

func WatchBinLog(conf config.Config) {
	q := New(config.Conf.RabbitMQ.RabbitURL, "queue.deleteCache")
	defer q.Close()

	q.Bind(config.Conf.RabbitMQ.CanalExchange, "")
	c := q.Consume()

	cacheDB := models.InitCacheDB(conf)
	conn := cacheDB.RedisPool.Get()

	for msg := range c {
		payload := getPayload(msg.Body)
		var err error
		// 监听到了mq发送过来的binlog变动，删除缓存
		_, err = conn.Do("HDEL", payload.UserIdentity, payload.ParentId+"file")
		_, err = conn.Do("HDEL", payload.UserIdentity, payload.ParentId+"folder")
		//_, err = conn.Do("DEL", payload.UserIdentity)
		logx.Info("Binlog有变动，删除缓存：", payload.UserIdentity, " ", payload.ParentId)

		// 如果失败，往mq中重新发送
		if err != nil {
			logx.Error("删除缓存失败, payload: ", payload)
			retryMq := New(config.Conf.RabbitMQ.RabbitURL, "")
			retryMq.Publish(config.Conf.RabbitMQ.CanalExchange, string(msg.Body))
		}
	}
}

func getPayload(data []byte) models.Payload {
	data1 := string(data)
	data2 := data1[strings.Index(data1, "user_identity")+16:]
	_userIdentity := data2[:(strings.Index(data2, ",") - 1)]

	data3 := data1[strings.Index(data1, "parent_id")+12:]
	_parentId := data3[:(strings.Index(data3, ",") - 1)]

	_userIdentity = strings.TrimRight(strings.Trim(_userIdentity, "\""), "\\")
	_parentId = strings.TrimRight(strings.Trim(_parentId, "\""), "\\")
	_userIdentity = strings.TrimLeft(strings.TrimLeft(_userIdentity, "\""), "\\")
	_parentId = strings.TrimLeft(strings.TrimLeft(_parentId, "\""), "\\")

	_userIdentity = _userIdentity[strings.Index(_userIdentity, "\\\"")+1:]
	_parentId = _parentId[strings.Index(_parentId, "\\\"")+1:]

	return models.Payload{
		UserIdentity: _userIdentity,
		ParentId:     _parentId,
	}
}
