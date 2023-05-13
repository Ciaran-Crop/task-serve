package config

import (
	"fmt"
)

const HOST string = "localhost"
const REDIS_PORT int = 6379
const RABBIT_PORT int = 5672
const RABBIT_MQ_NAME string = "task_mq"
const RABBIT_USER = "ciaran"
const RABBIT_PASSWORD = "123456"

type Task struct {
	TaskName    string
	TaskId      string
	TaskCommand string
}

type Status int

const (
	New Status = iota
	Ready
	Run
	Finish
	Error
)

func (t Task) String() string {
	return fmt.Sprintf("TaskName: %s, TaskId: %s, TaskCommand: %s\n", t.TaskName, t.TaskId, t.TaskCommand)
}
