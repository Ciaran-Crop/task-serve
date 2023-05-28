package rabbitConn

import (
	"log"
	"strconv"
	"task-serve/config"
	"task-serve/utils"

	"github.com/streadway/amqp"
)

type RabbitConnPool struct {
	Dial func() (*amqp.Connection, error)

	MinIdleConns int

	IdlesChan chan *amqp.Connection
}

var (
	Addr = "amqp://" + config.RABBIT_USER + ":" + config.RABBIT_PASSWORD + "@" + config.HOST + ":" + strconv.Itoa(config.RABBIT_PORT) + "/"
)

var rabbitConnPool *RabbitConnPool

func (rc *RabbitConnPool) InitPool() error {
	rc.IdlesChan = make(chan *amqp.Connection, rc.MinIdleConns)
	for i := 0; i < rc.MinIdleConns; i++ {
		rbmq, err := rc.Dial()
		if err != nil {
			return err
		}
		rc.IdlesChan <- rbmq
	}
	return nil
}

func (rc *RabbitConnPool) Get() *amqp.Connection {
	rbmq := <-rc.IdlesChan
	return rbmq
}

func (rc *RabbitConnPool) Release(rbmq *amqp.Connection) {
	rc.IdlesChan <- rbmq
}

func (rc *RabbitConnPool) ProduceTask(task *config.Task) error {
	rbmq := rc.Get()
	defer rc.Release(rbmq)
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

func (rc *RabbitConnPool) Consume() (*amqp.Channel, <-chan amqp.Delivery) {
	// 获取channel
	rbmq := rc.Get()
	defer rc.Release(rbmq)
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

func GetRabbitPool() *RabbitConnPool {
	return rabbitConnPool
}

func InitRabbitMQ() {
	rabbitConnPool = &RabbitConnPool{
		Dial: func() (*amqp.Connection, error) {
			return amqp.Dial(Addr)
		},
		MinIdleConns: 3,
	}

	rabbitConnPool.InitPool()
}

func CloseRabbitMQ() {
	rabbitConnPool = nil
}
