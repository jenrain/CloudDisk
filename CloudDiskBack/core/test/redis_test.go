package test

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"testing"
	"time"
)

// RedisPool 创建Redis全局连接池句柄
var RedisPool redis.Pool

// 初始化redis连接池
func init() {
	RedisPool = redis.Pool{
		MaxIdle:         20,
		MaxActive:       50,
		MaxConnLifetime: time.Duration(300),
		IdleTimeout:     time.Duration(60),
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				"116.62.177.68:63791",
				redis.DialPassword("501124524"),
			)
		},
	}
}

type cacheDb struct{}

var CacheDb = &cacheDb{}

// Set 写入数据的方法
func (c cacheDb) Set(key string, value interface{}) error {
	// 从连接池获取一条连接
	conn := RedisPool.Get()
	defer conn.Close()

	bytes, _ := json.Marshal(value)
	_, err := conn.Do("setex", key, 180, string(bytes))

	return err
}

// Get 获取数据的方法
func (c cacheDb) Get(key string, obj interface{}) error {
	// 从Redis连接池获取一个连接
	conn := RedisPool.Get()
	defer conn.Close()

	redisData, err := redis.Bytes(conn.Do("get", key))
	if err != nil {
		return err
	}
	// 将字节数组类型以json的形式写入obj
	json.Unmarshal(redisData, obj)
	return err
}

// Del 删除数据的方法
func (c cacheDb) Del(key string) (int64, error) {
	// 从Redis连接池获取一个连接
	conn := RedisPool.Get()
	defer conn.Close()

	redisData, err := conn.Do("del", key)
	count, _ := redisData.(int64)
	return count, err
}

func TestRedis(t *testing.T) {
	//err := CacheDb.Set("greet", "hello1")
	//if err != nil {
	//	t.Fatal("设置失败")
	//}
	//fmt.Println("存储成功")
	//var greet string
	//err := CacheDb.Get("hello", &greet)
	//if err != nil {
	//	t.Fatal("不存在")
	//}
	//fmt.Println("greet: ", greet)
	count, err := CacheDb.Del("greet")
	if err != nil {
		t.Fatal("待删除的键不存在")
	}
	fmt.Println("删除了", count, "个键")
}
