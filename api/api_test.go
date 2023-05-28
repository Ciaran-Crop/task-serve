package api_test

import (
	"fmt"
	"task-serve/api"
	"task-serve/config"
	"task-serve/rabbitConn"
	"task-serve/redisConn"
	"testing"
)

func TestCreateTask(t *testing.T) {
	rabbitConn.InitRabbitMQ()
	redisConn.InitRedis()
	taskId, err := api.CreateTask("test_create_task", "print('hello world')")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(taskId)
}

func TestSelectResult(t *testing.T) {
	rabbitConn.InitRabbitMQ()
	redisConn.InitRedis()
	taskId, err := api.CreateTask("test_create_task", "print('hello world')")
	if err != nil {
		t.Error(err.Error())
	}
	status, err := api.SelectResult(taskId)
	if err != nil {
		t.Error("select status error", err)
	}
	if status != "New" {
		t.Errorf("test select result expected status == Ready, but get status: %s", status)
	}
}

func TestUpdateState(t *testing.T) {
	rabbitConn.InitRabbitMQ()
	redisConn.InitRedis()
	taskId, err := api.CreateTask("test_create_task", "print('hello world')")
	if err != nil {
		t.Error(err.Error())
	}
	err = api.UpdateTaskStatus(taskId, config.Ready)
	if err != nil {
		t.Error(err.Error())
	}
	status, err := api.SelectResult(taskId)
	if err != nil {
		panic(err)
	}
	if status != "Ready" {
		t.Errorf("test select result expected status == Finish, but get status: %s", status)
	}
}
