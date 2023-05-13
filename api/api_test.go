package api_test

import (
	"fmt"
	"task-serve/api"
	"task-serve/config"
	"task-serve/rabbitConn"
	"task-serve/redisConn"
	"testing"
)

func InitConnection() {
	redisConn.InitRedis()
	rabbitConn.InitRabbitMQ()
}

func CloseConnection() {
	redisConn.CloseRedis()
	rabbitConn.CloseRabbitMQ()
}

func TestCreateTask(t *testing.T) {
	InitConnection()
	defer CloseConnection()
	taskId, err := api.CreateTask("test_create_task", "print('hello world')")
	if err != nil {
		t.Error("create task error", err)
	}
	fmt.Println(taskId)
}

func TestSelectResult(t *testing.T) {
	InitConnection()
	defer CloseConnection()
	taskId, err := api.CreateTask("test_create_task", "print('hello world')")
	if err != nil {
		t.Error("create task error", err)
	}
	status, err := api.SelectResult(taskId)
	if err != nil {
		t.Error("select status error", err)
	}
	if status != "Ready" {
		t.Errorf("test select result expected status == Ready, but get status: %s", status)
	}
}

func TestUpdateState(t *testing.T) {
	InitConnection()
	defer CloseConnection()
	taskId, err := api.CreateTask("test_create_task", "print('hello world')")
	if err != nil {
		t.Error("create task error", err)
	}
	err = api.UpdateTaskStatus(taskId, config.Finish)
	if err != nil {
		t.Error("create task error", err)
	}
	status, err := api.SelectResult(taskId)
	if err != nil {
		panic(err)
	}
	if status != "Finish" {
		t.Errorf("test select result expected status == Finish, but get status: %s", status)
	}
}
