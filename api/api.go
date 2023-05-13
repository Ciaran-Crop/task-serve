package api

import (
	"fmt"
	"strconv"
	"task-serve/config"
	"task-serve/rabbitConn"
	"task-serve/redisConn"
	"time"
)

func CreateTask(taskName string, taskCommand string) (string, error) {
	global_task_id, err := redisConn.RedisIncr("global_task_id")
	if err != nil {
		return "", err
	}
	taskId := "task-" + strconv.Itoa(global_task_id)
	task := config.Task{
		TaskName:    taskName,
		TaskId:      taskId,
		TaskCommand: taskCommand,
	}
	err = redisConn.RedisSet(taskId, int(config.New), time.Hour*2)
	if err != nil {
		return "", err
	}
	err = rabbitConn.ProduceTask(task)
	if err != nil {
		redisConn.RedisDel(taskId)
	}
	err = redisConn.RedisSet(taskId, int(config.Ready), time.Hour*2)
	if err != nil {
		redisConn.RedisDel(taskId)
		return "", err
	}
	fmt.Printf("Create Task : %s", task)
	return taskId, nil
}

func decodeStatus(status int) string {
	var result string
	switch status {
	case int(config.New):
		result = "New"
	case int(config.Ready):
		result = "Ready"
	case int(config.Run):
		result = "Run"
	case int(config.Finish):
		result = "Finish"
	default:
		result = "Error"
	}
	return result
}

func SelectResult(taskId string) (string, error) {
	status, err := redisConn.RedisGet(taskId)
	if err != nil {
		return "", err
	}
	val, err := strconv.Atoi(status)
	if err != nil {
		return "", err
	}
	result := decodeStatus(val)
	return result, nil
}

func UpdateTaskStatus(taskId string, status config.Status) error {
	err := redisConn.RedisSet(taskId, int(status), time.Hour*2)
	if err != nil {
		return err
	}
	fmt.Printf("Update Task: %s to %s\n", taskId, decodeStatus(int(status)))
	return nil
}
