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
	redisClient, redisCtx, err := redisConn.GetClient()
	rabbitConnPool := rabbitConn.GetRabbitPool()
	if err != nil {
		return "", err
	}
	val, err := redisClient.IncrBy(redisCtx, "global_task_id", 1).Result()
	if err != nil {
		return "", err
	}
	global_task_id := int(val)
	taskId := "task-" + strconv.Itoa(global_task_id)
	now_time := time.Now().UnixMilli()
	task := &config.Task{
		TaskName:    taskName,
		TaskId:      taskId,
		TaskCommand: taskCommand,
		TaskTime:    now_time,
		TaskStatus:  config.New,
	}
	_, err = redisClient.HSet(redisCtx, taskId, task.GetMap()).Result()
	if err != nil {
		return taskId, err
	}

	err = rabbitConnPool.ProduceTask(task)
	if err != nil {
		redisClient.HSet(redisCtx, taskId, "Status", config.Error)
		return taskId, err
	}

	return taskId, nil
}

func CancelTask(taskId string) error {
	rabbitConnPool := rabbitConn.GetRabbitPool()
	now_time := time.Now().UnixMilli()
	task := &config.Task{
		TaskName:    "cancel-task",
		TaskId:      taskId,
		TaskCommand: "cancel",
		TaskTime:    now_time,
		TaskStatus:  config.Cancel,
	}
	err := rabbitConnPool.ProduceTask(task)
	if err != nil {
		return err
	}
	return nil
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
	case int(config.Cancel):
		result = "Cancel"
	default:
		result = "Error"
	}
	return result
}

func SelectResult(taskId string) (string, error) {
	redisClient, redisCtx, err := redisConn.GetClient()
	if err != nil {
		return "", err
	}
	status, err := redisClient.HGet(redisCtx, taskId, "Status").Result()
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
	redisClient, redisCtx, err := redisConn.GetClient()
	if err != nil {
		return err
	}
	status1, err := redisClient.HGet(redisCtx, taskId, "Status").Result()
	if err != nil {
		return err
	}
	beforeStatus, err := strconv.Atoi(status1)
	if err != nil {
		return err
	}
	if !checkUpdate(config.Status(beforeStatus), status) {
		if config.Status(beforeStatus) == config.Cancel {
			return nil
		}
		return fmt.Errorf("can't update status from %v to %v", decodeStatus(beforeStatus), decodeStatus(int(status)))
	}
	_, err = redisClient.HSet(redisCtx, taskId, "Status", status).Result()
	if err != nil {
		return err
	}
	fmt.Printf("Update Task: %s to %s\n", taskId, decodeStatus(int(status)))
	return nil
}

// New -> Ready -> Run -> Finish
// ALL(except Finish) -> Error
// New,Ready,Run -> Cancel
func checkUpdate(status1 config.Status, status2 config.Status) bool {
	if status2 == config.Error && status1 != config.Finish {
		return true
	}
	switch status1 {
	case config.New:
		if status2 == config.Ready || status2 == config.Cancel {
			return true
		}
	case config.Ready:
		if status2 == config.Run || status2 == config.Cancel {
			return true
		}
	case config.Run:
		if status2 == config.Finish || status2 == config.Cancel {
			return true
		}
	}
	return false
}
