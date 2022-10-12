package models

import (
	"core/internal/config"
	"github.com/gomodule/redigo/redis"
	"time"
)

// RedisPool 创建Redis全局连接池句柄
var redisPool redis.Pool

type CacheDb struct {
	RedisPool redis.Pool
}

var cacheDB CacheDb

// InitCacheDB 初始化redis连接池
func InitCacheDB(c config.Config) CacheDb {
	redisPool = redis.Pool{
		MaxIdle:         20,
		MaxActive:       50,
		MaxConnLifetime: time.Duration(300),
		IdleTimeout:     time.Duration(60),
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				c.Redis.Addr,
				redis.DialPassword(c.Redis.Password),
			)
		},
	}
	cacheDB.RedisPool = redisPool
	return cacheDB
}

// Set 写入数据的方法
//func (c CacheDb) Set(key string, value interface{}) error {
//	// 从连接池获取一条连接
//	conn := RedisPool.Get()
//	defer conn.Close()
//
//	bytes, _ := json.Marshal(value)
//	_, err := conn.Do("setex", key, 3600, string(bytes))
//
//	return err
//}
//
//// Get 获取数据的方法
//func (c CacheDb) Get(key string, obj interface{}) error {
//	// 从Redis连接池获取一个连接
//	conn := RedisPool.Get()
//	defer conn.Close()
//
//	redisData, err := redis.Bytes(conn.Do("get", key))
//	if err != nil {
//		return err
//	}
//	// 将字节数组类型以json的形式写入obj
//	json.Unmarshal(redisData, obj)
//	return err
//}
//
//// Del 删除数据的方法
//func (c CacheDb) Del(key string) (int64, error) {
//	// 从Redis连接池获取一个连接
//	conn := RedisPool.Get()
//	defer conn.Close()
//
//	redisData, err := conn.Do("del", key)
//	count, _ := redisData.(int64)
//	return count, err
//}
