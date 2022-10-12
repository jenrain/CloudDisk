package test

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"testing"
)

// RabbitMQ RabbitMQ结构图
type rabbitMQ struct {
	channel  *amqp.Channel
	Name     string
	exchange string
}

// New 连接RabbitMQ服务，声明一个消息队列
func new1(s string) *rabbitMQ {
	conn, e := amqp.Dial(s)
	if e != nil {
		panic(e)
	}

	ch, e := conn.Channel()
	if e != nil {
		panic(e)
	}

	q, e := ch.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil)
	if e != nil {
		panic(e)
	}

	mq := new(rabbitMQ)
	mq.channel = ch
	mq.Name = q.Name
	return mq
}

// Bind 消息队列绑定交换机
func (q *rabbitMQ) bind(exchange, key string) {
	e := q.channel.QueueBind(
		q.Name,
		key,
		exchange,
		false,
		nil)
	if e != nil {
		panic(e)
	}

	q.exchange = exchange
}

// Send 向消息队列发布消息
func (q *rabbitMQ) send(queue string, body interface{}) {
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}

	e = q.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		panic(e)
	}
}

// Publish 向交换机发送消息
func (q *rabbitMQ) publish(excahnge string, body interface{}) {
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}

	e = q.channel.Publish(
		excahnge,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	if e != nil {
		panic(e)
	}
}

func (q *rabbitMQ) close() {
	q.channel.Close()
}

func TestMQ(t *testing.T) {
	//server := mq.ChooseRandomDataServer()
	//fmt.Println(server.HuaweiObsBucketRootFolder)
	//servers := mq.GetDataServers()
	//fmt.Println("map中的数据：", servers)
	q := new1("amqp://guest:guest@116.62.177.68:5672/")
	defer q.close()

	q.bind("canal.deleteCache", "")

	q.publish("canal.deleteCache", "hell2")
}
