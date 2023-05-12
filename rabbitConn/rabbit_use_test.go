package rabbitConn_test

import (
	"fmt"
	"task-serve/config"
	"task-serve/rabbitConn"
	"testing"
)

func TestConn(t *testing.T) {
	err := rabbitConn.InitRabbitMQ()
	if err != nil {
		t.Errorf("Connect Error!")
		fmt.Println(err)
	}
}

func TestProduce(t *testing.T) {
	rabbitConn.InitRabbitMQ()
	err := rabbitConn.ProduceTask(config.Task{
		TaskName:    "Test",
		TaskId:      "0",
		TaskCommand: "print('hello world')",
	})
	if err != nil {
		t.Errorf("Produce Error!")
		fmt.Println(err)
	}
}

func TestConsume(t *testing.T) {
	rabbitConn.InitRabbitMQ()
	task := rabbitConn.Consume()
	fmt.Println(task.TaskId)
}
