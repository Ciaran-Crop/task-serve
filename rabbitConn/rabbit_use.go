package rabbitConn

import (
	"log"
	"strconv"
	"task-serve/config"
	"task-serve/utils"

	"github.com/streadway/amqp"
)

var rbmq *amqp.Connection
var err error

func InitRabbitMQ() error {
	Addr := "amqp://" + config.RABBIT_USER + ":" + config.RABBIT_PASSWORD + "@" + config.HOST + ":" + strconv.Itoa(config.RABBIT_PORT) + "/"
	if rbmq == nil {
		rbmq, err = amqp.Dial(Addr)
		if err != nil {
			return err
		}
	}
	return nil
}

func CloseRabbitMQ() {
	if rbmq != nil {
		rbmq.Close()
	}
}

func ProduceTask(task config.Task) error {
	// 获取channel
	ch, err := rbmq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	// 创建队列
	q, err := ch.QueueDeclare(
		config.RABBIT_MQ_NAME,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	// 发送消息
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(utils.Encode(task)),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func Consume() (*amqp.Channel, <-chan amqp.Delivery) {
	// 获取channel
	ch, err := rbmq.Channel()
	if err != nil {
		log.Fatal(err)
	}
	// 创建队列
	q, err := ch.QueueDeclare(
		config.RABBIT_MQ_NAME,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	//设置每次从消息队列获取任务的数量
	err = ch.Qos(
		1,     //预取任务数量
		0,     //预取大小
		false, //全局设置
	)

	if err != nil {
		log.Fatal(err)
	}
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	return ch, msgs
}
