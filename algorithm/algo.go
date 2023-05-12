package algorithm

import (
	"fmt"
	"task-serve/api"
	"task-serve/config"
	"task-serve/rabbitConn"
)

func Algo(task config.Task) {
	//
	api.UpdateTaskStatus(task.TaskId, config.Finish)
}

func RunServe() {
	fmt.Println("Start Algorithm Serve")
	go func() {
		task := rabbitConn.Consume()
		api.UpdateTaskStatus(task.TaskId, config.Run)
		go Algo(task)
	}()
}
