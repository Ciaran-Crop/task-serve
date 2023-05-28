package main

import (
	"fmt"
	"sync"
	"task-serve/algorithm"
	"task-serve/api"
	"task-serve/config"
	"task-serve/rabbitConn"
	"task-serve/redisConn"
	"testing"
	"time"
)

func TestOneTask(t *testing.T) {
	rabbitConn.InitRabbitMQ()
	redisConn.InitRedis()
	go algorithm.RunServe()
	taskId, err := api.CreateTask("test_task", "print('Hello World!')")
	if err != nil {
		t.Errorf(err.Error())
	}
	for {
		status, err := api.SelectResult(taskId)
		if status == "Finish" {
			break
		}
		time.Sleep(time.Second)
		if err != nil {
			t.Errorf(err.Error())
		}
		fmt.Printf("Now Task Status: %s\n", status)
	}
}

func TestMultiTask(t *testing.T) {
	rabbitConn.InitRabbitMQ()
	redisConn.InitRedis()
	go algorithm.RunServe()
	var wg sync.WaitGroup
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func() {
			taskId, err := api.CreateTask("test_task", "print('Hello World!')")
			if err != nil {
				t.Errorf(err.Error())
			}
			for {
				status, err := api.SelectResult(taskId)
				if status == "Finish" {
					break
				}
				time.Sleep(time.Second)
				if err != nil {
					t.Errorf(err.Error())
				}
				fmt.Printf("Now Task Status: %s\n", status)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestCancelOneTask(t *testing.T) {
	rabbitConn.InitRabbitMQ()
	redisConn.InitRedis()
	go algorithm.RunServe()
	taskId, err := api.CreateTask("test_task", "print('Hello World!')")
	if err != nil {
		t.Errorf(err.Error())
	}
	for {
		api.UpdateTaskStatus(taskId, config.Cancel)
		status, err := api.SelectResult(taskId)
		if status == "Finish" || status == "Cancel" {
			break
		}
		time.Sleep(time.Second)
		if err != nil {
			t.Errorf(err.Error())
		}
		fmt.Printf("Now Task Status: %s\n", status)
	}
}

func TestErrorOneTask(t *testing.T) {
	rabbitConn.InitRabbitMQ()
	redisConn.InitRedis()
	go algorithm.RunServe()
	taskId, err := api.CreateTask("test_task", "print('Hello World!')")
	if err != nil {
		t.Errorf(err.Error())
	}
	err = api.UpdateTaskStatus(taskId, config.Run)
	if err != nil {
		fmt.Print(err.Error())
	}
}

func TestCancelMultiTask(t *testing.T) {
	rabbitConn.InitRabbitMQ()
	redisConn.InitRedis()
	go algorithm.RunServe()
	var wg sync.WaitGroup
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func() {
			taskId, err := api.CreateTask("test_task", "print('Hello World!')")
			if err != nil {
				t.Errorf(err.Error())
			}
			for {
				api.UpdateTaskStatus(taskId, config.Cancel)
				status, err := api.SelectResult(taskId)
				if status == "Finish" || status == "Cancel" {
					break
				}
				time.Sleep(time.Second)
				if err != nil {
					t.Errorf(err.Error())
				}
				fmt.Printf("Now Task Status: %s\n", status)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
