package api

import (
	"strconv"
	"task-serve/config"
	"task-serve/redisConn"
	"time"
)

func CreateTask(taskName string, taskCommand string) string {
	return ""
}

func SelectResult(taskId string) string {
	status, err := redisConn.RedisGet(taskId)
	if err != nil {
		panic(err)
	}
	var result string
	val, err := strconv.Atoi(status)
	if err != nil {
		panic(err)
	}
	switch val {
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

func UpdateTaskStatus(taskId string, status config.Status) {
	err := redisConn.RedisSet(taskId, status, time.Hour*24)
	if err != nil {
		panic(err)
	}
}
