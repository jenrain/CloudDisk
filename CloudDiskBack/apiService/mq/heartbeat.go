package mq

import (
	"core/internal/config"
	"core/models"
	"encoding/json"
	"math"
	"math/rand"
	"sync"
	"time"
)

// 数据服务，键是服务地址，值是最后心跳时间
var dataServers sync.Map

// ListenHeartbeat 监听数据服务心跳
func ListenHeartbeat() {
	q := New(config.Conf.RabbitMQ.RabbitURL, "")
	defer q.Close()

	q.Bind(config.Conf.RabbitMQ.ApiServers, "")
	c := q.Consume()

	go removeExpireDataServer()

	for msg := range c {
		var o models.Obs
		json.Unmarshal(msg.Body, &o)
		//fmt.Println("接收到的心跳节点：", o.HuaweiObsBucket)
		dataServers.Store(o, time.Now())
	}
}

// 移除过期的数据服务
func removeExpireDataServer() {
	for {
		dataServers.Range(func(key, value interface{}) bool {
			t, _ := value.(time.Time)
			if t.Add(10 * time.Second).Before(time.Now()) {
				dataServers.Delete(key)
			}
			return true
		})
	}
}

func GetDataServers() []models.Obs {
	ds := make([]models.Obs, 0)

	dataServers.Range(func(key, value interface{}) bool {
		obs, _ := key.(models.Obs)
		ds = append(ds, obs)
		return true
	})

	return ds
}

//ChooseRandomDataServers 返回n个随机数据服务地址
func ChooseRandomDataServers() (servers []models.Obs) {
	ds := GetDataServers()
	n := len(ds)

	if n == 0 {
		return nil
	}
	randNumbers := rand.Perm(n)
	count := int(math.Ceil(float64(n) / 3))
	for i := 0; i < count; i++ {
		servers = append(servers, ds[randNumbers[i]])
	}
	return
}

// ChooseRandomDataServer 返回一个数据服务地址
func ChooseRandomDataServer() (server models.Obs) {
	ds := GetDataServers()
	n := len(ds)

	if n == 0 {
		return
	}
	return ds[rand.Intn(n)]
}

//ChooseRandomDataServersExcept 返回除了已经选择过的数据服务地址
func ChooseRandomDataServersExcept(node models.Obs) (servers []models.Obs) {
	ds := GetDataServers()
	n := len(ds)

	if n == 0 {
		return nil
	}
	randNumbers := rand.Perm(n)
	count := int(math.Ceil(float64(n)/3)) - 1
	j := -1
	for {
		if len(servers) == count {
			break
		}
		j++
		if ds[randNumbers[j]].HuaweiObsBucket == node.HuaweiObsBucket {
			continue
		}
		servers = append(servers, ds[randNumbers[j]])
	}
	return
}

// GetNodeCount 获取当前节点的数量
func getNodeCount() (count int) {
	dataServers.Range(func(key, value interface{}) bool {
		if key != "" && key != " " {
			count++
		}
		return true
	})
	return
}
