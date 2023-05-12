package redisConn_test

import (
	"strconv"
	"task-serve/config"
	"task-serve/redisConn"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	redisConn.InitRedis()
	if ans, err := redisConn.RedisGet("test_a"); ans != "bbb" || err != nil {
		t.Errorf("test_a expected be 'bbb', but get %s", ans)
	}
}

func TestSet(t *testing.T) {
	redisConn.InitRedis()
	err := redisConn.RedisSet("test_b", "bbb", time.Hour)
	if err != nil {
		t.Error("redis set error", err)
	}
	if ans, err := redisConn.RedisGet("test_b"); ans != "bbb" || err != nil {
		t.Errorf("test_a expected be 'bbb', but get %s", ans)
	}
}

func TestSetStatus(t *testing.T) {
	redisConn.InitRedis()
	err := redisConn.RedisSet("test_c", int(config.New), time.Hour)
	if err != nil {
		t.Error("redis set error", err)
	}
	ans, err := redisConn.RedisGet("test_c")
	if err != nil {
		t.Error("redis get error", err)
	}
	if status, err := strconv.Atoi(ans); status != int(config.New) || err != nil {
		t.Errorf("test_a expected be '0', but get %s", ans)
	}
}
