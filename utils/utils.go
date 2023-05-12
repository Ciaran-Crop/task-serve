package utils

import (
	"bytes"
	"encoding/gob"
	"task-serve/config"
)

func Encode(data interface{}) []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(data)
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func Decode(data []byte) config.Task {
	b := bytes.NewBuffer(data)
	var task config.Task
	dec := gob.NewDecoder(b)
	err := dec.Decode(&task)
	if err != nil {
		panic(err)
	}
	return task
}
