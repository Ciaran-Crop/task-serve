package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

func TestHttpCreate(t *testing.T) {
	url := "http://localhost:9001/create_task"
	rep, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader("task_name=http_create_test&task_command=HelloWorld"))
	if err != nil {
		t.Error(err.Error())
	}
	defer rep.Body.Close()
	bytes, _ := ioutil.ReadAll(rep.Body)
	m := make(map[string]interface{})
	json.Unmarshal(bytes, &m)
	fmt.Println(m)
	fmt.Println(m["result"], m["task_id"])
}

func TestHttpGetOneTask(t *testing.T) {
	taskId := "task-43"
	// rep, err := http.Post("http://localhost:9001/view_task", "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("task_id=%s", taskId)))
	rep, err := http.Get("http://localhost:9001/view_task" + "?task_id=" + taskId)
	if err != nil {
		t.Error(err.Error())
	}
	defer rep.Body.Close()
	bytes, _ := ioutil.ReadAll(rep.Body)
	m := make(map[string]interface{})
	json.Unmarshal(bytes, &m)
	fmt.Println(m)
}

func TestHttpGetTaskList(t *testing.T) {
	rep, err := http.Get("http://localhost:9001/get_task_list")
	if err != nil {
		t.Error(err.Error())
	}
	defer rep.Body.Close()
	bytes, _ := ioutil.ReadAll(rep.Body)
	m := make(map[string]interface{})
	json.Unmarshal(bytes, &m)
	fmt.Println(m)
}

func TestHttpCancel(t *testing.T) {
	url := "http://localhost:9001/create_task"
	rep, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader("task_name=http_create_test&task_command=HelloWorld"))
	if err != nil {
		t.Error(err.Error())
	}
	defer rep.Body.Close()
	bytes, _ := ioutil.ReadAll(rep.Body)
	m := make(map[string]interface{})
	json.Unmarshal(bytes, &m)
	fmt.Println(m)
	// time.Sleep(time.Second * 1)
	rep, err = http.Get(fmt.Sprintf("http://localhost:9001/cancel_task?task_id=%s", m["task_id"]))
	if err != nil {
		t.Error(err.Error())
	}
	bytes, _ = ioutil.ReadAll(rep.Body)
	m = make(map[string]interface{})
	json.Unmarshal(bytes, &m)
	fmt.Println(m)
}
