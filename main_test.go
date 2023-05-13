package main

import (
	"fmt"
	"sync"
	"task-serve/algorithm"
	"task-serve/api"
	"testing"
	"time"
)

func errorPrint(t *testing.T, err error, msg string) {
	if err != nil {
		t.Error(msg, err)
	}
}

func TestOneTask(t *testing.T) {
	InitConnection()
	defer CloseConnection()
	go algorithm.RunServe()
	taskId, err := api.CreateTask("test_task", "print('Hello World!')")
	errorPrint(t, err, "create task error")
	for {
		status, err := api.SelectResult(taskId)
		if status == "Finish" {
			break
		}
		time.Sleep(time.Second)
		errorPrint(t, err, "select result error")
		fmt.Printf("Now Task Status: %s\n", status)
	}
}

func TestMultiTask(t *testing.T) {
	InitConnection()
	defer CloseConnection()
	go algorithm.RunServe()
	var wg sync.WaitGroup
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func() {
			taskId, err := api.CreateTask("test_task", "print('Hello World!')")
			errorPrint(t, err, "create task error")
			for {
				status, err := api.SelectResult(taskId)
				if status == "Finish" {
					break
				}
				time.Sleep(time.Second)
				errorPrint(t, err, "select result error")
				fmt.Printf("Now Task Status: %s\n", status)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
