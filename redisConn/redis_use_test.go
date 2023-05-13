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
		t.Errorf("test_get expected be 'bbb', but get %s", ans)
	}
}

func TestSet(t *testing.T) {
	redisConn.InitRedis()
	err := redisConn.RedisSet("test_b", "bbb", time.Hour)
	if err != nil {
		t.Error("redis set error", err)
	}
	if ans, err := redisConn.RedisGet("test_b"); ans != "bbb" || err != nil {
		t.Errorf("test_set expected be 'bbb', but get %s", ans)
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
		t.Errorf("test_status expected be '0', but get %s", ans)
	}
}

func TestIncr(t *testing.T) {
	redisConn.InitRedis()
	oldVal, err := redisConn.RedisGet("test_d")
	if err != nil {
		panic(err)
	}
	val, err := redisConn.RedisIncr("test_d")
	if err != nil {
		t.Error("redis incr error", err)
	}
	oldIntVal, err := strconv.Atoi(oldVal)
	if err != nil {
		panic(err)
	}
	if oldIntVal != val-1 {
		t.Errorf("test_incr expected be oldIntVal == val - 1, but get oldIntVal = %d, val = %d", oldIntVal, val)
	}
}

func TestDel(t *testing.T) {
	redisConn.InitRedis()
	redisConn.RedisDel("test_e")
	if val, err := redisConn.RedisGet("test_e"); err == nil {
		t.Errorf("test_del expected get errror, but get val %s", val)
	}
}
