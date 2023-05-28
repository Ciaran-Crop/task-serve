package algorithm

import (
	"context"
	"fmt"
	"task-serve/api"
	"task-serve/config"
	"task-serve/rabbitConn"
	"task-serve/taskOp"
	"task-serve/utils"
)

var wait chan bool
var ctxMap map[string]context.CancelFunc

func doAlgo(task *config.Task, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := api.UpdateTaskStatus(task.TaskId, config.Run)
			taskFail(task.TaskId, err)
			fmt.Printf("Start Task : %s", task)
			taskOp.DoHelloWorld()
			err = api.UpdateTaskStatus(task.TaskId, config.Finish)
			taskFail(task.TaskId, err)
			delete(ctxMap, task.TaskId)
			return
		}
	}
}

func taskFail(taskId string, err error) {
	if err != nil {
		api.UpdateTaskStatus(taskId, config.Error)
	}
}

func RunServe() {
	fmt.Println("Start Algorithm Serve")
	ctxMap = make(map[string]context.CancelFunc)
	rabbitConnPool := rabbitConn.GetRabbitPool()
	ch, msgs := rabbitConnPool.Consume()
	defer ch.Close()
	go func() {
		for v := range msgs {
			task := utils.Decode(v.Body)
			v.Ack(true)
			status, err := api.SelectResult(task.TaskId)
			taskFail(task.TaskId, err)
			if status == "New" {
				fmt.Printf("Get Task : %s", task)
				err := api.UpdateTaskStatus(task.TaskId, config.Ready)
				taskFail(task.TaskId, err)
				ctx, cancel := context.WithCancel(context.Background())
				ctxMap[task.TaskId] = cancel
				go doAlgo(task, ctx)
			} else if status == "Cancel" {
				CancelTask(task.TaskId)
			}
		}
	}()
	<-wait
}

func CancelTask(taskId string) error {
	if cancel, ok := ctxMap[taskId]; ok {
		cancel()
	}
	err := api.UpdateTaskStatus(taskId, config.Cancel)
	if err != nil {
		return err
	}
	return nil
}
