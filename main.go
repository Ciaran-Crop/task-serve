package main

import (
	"flag"
	"fmt"
	"task-serve/algorithm"
	"task-serve/api"
	"task-serve/rabbitConn"
	"task-serve/redisConn"
)

func InitConnection() {
	redisConn.InitRedis()
	rabbitConn.InitRabbitMQ()
}

func CloseConnection() {
	redisConn.CloseRedis()
	rabbitConn.CloseRabbitMQ()
}

func main() {
	var serve = flag.String("s", "api", "选择启动服务[api / algo]")
	var op = flag.String("op", "select", "选择操作[select / create]")
	var taskName = flag.String("tname", "", "任务名")
	var taskId = flag.String("tid", "", "任务id")
	var taskCommand = flag.String("tcommand", "", "任务命令")
	InitConnection()
	defer CloseConnection()
	flag.Parse()
	if *serve == "api" {
		if *op == "select" {
			if *taskId == "" {
				flag.Usage()
				return
			}
			status, err := api.SelectResult(*taskId)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Task id is : %s, status: %s", *taskId, status)
		} else if *op == "create" {
			if *taskName == "" || *taskCommand == "" {
				flag.Usage()
				return
			}
			taskId, err := api.CreateTask(*taskName, *taskCommand)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Your Task id is : %s", taskId)
			return
		} else {
			flag.Usage()
			return
		}
	} else if *serve == "algo" {
		algorithm.RunServe()
	} else {
		flag.Usage()
	}
}
