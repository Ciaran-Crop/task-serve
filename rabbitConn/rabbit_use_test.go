package rabbitConn_test

import (
	"fmt"
	"task-serve/config"
	"task-serve/rabbitConn"
	"task-serve/redisConn"
	"task-serve/utils"
	"testing"
)

func TestConn(t *testing.T) {
	err := rabbitConn.InitRabbitMQ()
	if err != nil {
		t.Errorf("Connect Error!")
		fmt.Println(err)
	}
}

func InitConnection() {
	redisConn.InitRedis()
	rabbitConn.InitRabbitMQ()
}

func CloseConnection() {
	redisConn.CloseRedis()
	rabbitConn.CloseRabbitMQ()
}

func TestProduce(t *testing.T) {
	InitConnection()
	defer CloseConnection()
	err := rabbitConn.ProduceTask(config.Task{
		TaskName:    "Test",
		TaskId:      "2",
		TaskCommand: "print('hello world')",
	})
	if err != nil {
		t.Errorf("Produce Error!")
		fmt.Println(err)
	}
}

func TestConsume(t *testing.T) {
	InitConnection()
	defer CloseConnection()
	ch, msgs := rabbitConn.Consume()
	defer ch.Close()
	for v := range msgs {
		task := utils.Decode(v.Body)
		v.Ack(true)
		fmt.Printf("Get Task : %s", task)
	}
}
