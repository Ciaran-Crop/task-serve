package algorithm

import (
	"fmt"
	"task-serve/api"
	"task-serve/config"
	"task-serve/rabbitConn"
	"task-serve/taskOp"
	"task-serve/utils"
)

var wait chan bool

func doAlgo(task config.Task) {
	fmt.Printf("Start Task : %s", task)
	taskOp.DoHelloWorld()
	api.UpdateTaskStatus(task.TaskId, config.Finish)
}

func RunServe() {
	fmt.Println("Start Algorithm Serve")
	ch, msgs := rabbitConn.Consume()
	defer ch.Close()
	go func() {
		for v := range msgs {
			task := utils.Decode(v.Body)
			v.Ack(true)
			fmt.Printf("Get Task : %s", task)
			api.UpdateTaskStatus(task.TaskId, config.Run)
			go doAlgo(task)
		}
	}()
	<-wait
}
