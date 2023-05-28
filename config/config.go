package config

import (
	"encoding/json"
	"fmt"
)

const HOST string = "localhost"
const REDIS_PORT int = 6379
const RABBIT_PORT int = 5672
const RABBIT_MQ_NAME string = "task_mq"
const RABBIT_USER = "ciaran"
const RABBIT_PASSWORD = "123456"

type Status int

type Task struct {
	TaskName    string
	TaskId      string
	TaskCommand string
	TaskTime    int64
	TaskStatus  Status
}

const (
	New Status = iota
	Ready
	Run
	Finish
	Cancel
	Error
)

func (s Status) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s Status) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}

func (s *Task) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Task) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}

func (t *Task) String() string {
	return fmt.Sprintf("TaskName: %v, TaskId: %v, TaskCommand: %v, TaskTime: %v, TaskStatus: %v\n", t.TaskName, t.TaskId, t.TaskCommand, t.TaskTime, t.TaskStatus)
}

func (t *Task) GetMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["TaskName"] = t.TaskName
	m["TaskId"] = t.TaskId
	m["TaskCommand"] = t.TaskCommand
	m["TaskTime"] = t.TaskTime
	m["Status"] = t.TaskStatus
	return m
}
