package rabbitConn_test

import (
	"fmt"
	"task-serve/config"
	"task-serve/rabbitConn"
	"task-serve/utils"
	"testing"
	"time"
)

func TestProduce(t *testing.T) {
	rabbitConn.InitRabbitMQ()
	defer rabbitConn.CloseRabbitMQ()
	rabbitConnPool := rabbitConn.GetRabbitPool()
	err := rabbitConnPool.ProduceTask(&config.Task{
		TaskName:    "Test",
		TaskId:      "2",
		TaskCommand: "print('hello world')",
		TaskTime:    time.Now().UnixMilli(),
		TaskStatus:  config.New,
	})
	if err != nil {
		t.Errorf("Produce Error!")
		fmt.Println(err)
	}
}

func TestConsume(t *testing.T) {
	rabbitConn.InitRabbitMQ()
	defer rabbitConn.CloseRabbitMQ()
	rabbitConnPool := rabbitConn.GetRabbitPool()
	ch, msgs := rabbitConnPool.Consume()
	defer ch.Close()
	for v := range msgs {
		task := utils.Decode(v.Body)
		v.Ack(true)
		fmt.Printf("Get Task : %v", task)
	}
}
