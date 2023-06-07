package redisConn_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/Ciaran-crop/task-serve/config"
	"github.com/Ciaran-crop/task-serve/redisConn"
)

func TestGet(t *testing.T) {
	redisConn.InitRedis()
	redisClient, redisCtx, err := redisConn.GetClient()
	if err != nil {
		t.Errorf(err.Error())
	}
	if ans, err := redisClient.Get(redisCtx, "test_a").Result(); ans != "bbb" || err != nil {
		t.Errorf("test_get expected be 'bbb', but get %s", ans)
	}
}

func TestSet(t *testing.T) {
	redisConn.InitRedis()
	redisClient, redisCtx, err := redisConn.GetClient()
	if err != nil {
		t.Errorf(err.Error())
	}
	if _, err := redisClient.Set(redisCtx, "test_b", "bbb", 0).Result(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestTask(t *testing.T) {
	redisConn.InitRedis()
	task := &config.Task{
		TaskName:    "test-task",
		TaskId:      "test_c",
		TaskCommand: "none",
		TaskTime:    time.Now().UnixMilli(),
		TaskStatus:  config.Cancel,
	}
	redisClient, redisCtx, err := redisConn.GetClient()
	if err != nil {
		t.Errorf(err.Error())
	}
	if _, err := redisClient.HSet(redisCtx, task.TaskId, task.GetMap()).Result(); err != nil {
		t.Errorf(err.Error())
	}
	val, err := redisClient.HGet(redisCtx, task.TaskId, "Status").Result()
	if err != nil {
		t.Errorf(err.Error())
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		t.Errorf(err.Error())
	}
	if intVal != int(config.Cancel) {
		t.Errorf("HSet status failed : %v", intVal)
	}
	if _, err = redisClient.HSet(redisCtx, task.TaskId, "Status", config.Finish).Result(); err != nil {
		t.Errorf(err.Error())
	}
	val, err = redisClient.HGet(redisCtx, task.TaskId, "Status").Result()
	if err != nil {
		t.Errorf(err.Error())
	}
	intVal, err = strconv.Atoi(val)
	if err != nil {
		t.Errorf(err.Error())
	}
	if intVal != int(config.Finish) {
		t.Errorf("HSet status failed : %v", intVal)
	}
}

func TestIncr(t *testing.T) {
	redisConn.InitRedis()
	redisClient, redisCtx, err := redisConn.GetClient()
	if err != nil {
		t.Errorf(err.Error())
	}
	oldVal, err := redisClient.Get(redisCtx, "test_d").Result()
	if err != nil {
		t.Errorf(err.Error())
	}
	val, err := redisClient.IncrBy(redisCtx, "test_d", 1).Result()
	if err != nil {
		t.Errorf(err.Error())
	}
	oldIntVal, err := strconv.Atoi(oldVal)
	if err != nil {
		t.Errorf(err.Error())
	}
	if oldIntVal != int(val-1) {
		t.Errorf("test_incr expected be oldIntVal == val - 1, but get oldIntVal = %d, val = %d", oldIntVal, val)
	}
}

func TestDel(t *testing.T) {
	redisConn.InitRedis()
	redisClient, redisCtx, err := redisConn.GetClient()
	if err != nil {
		t.Errorf(err.Error())
	}
	if _, err := redisClient.Del(redisCtx, "test_e").Result(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestDelAllTask(t *testing.T) {
	redisConn.InitRedis()
	redisClient, redisCtx, err := redisConn.GetClient()
	if err != nil {
		t.Errorf(err.Error())
	}
	taskList, _ := redisClient.Keys(redisCtx, "task-*").Result()
	redisClient.Del(redisCtx, taskList...)
}
