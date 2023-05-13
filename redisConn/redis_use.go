package redisConn

import (
	"context"
	"fmt"
	"strconv"
	"task-serve/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx context.Context

func InitRedis() {
	Addr := config.HOST + ":" + strconv.Itoa(config.REDIS_PORT)
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr: Addr,
			DB:   0,
		})
	}
	if ctx == nil {
		ctx = context.Background()
	}
}

func CloseRedis() {
	if rdb != nil {
		rdb.Close()
	}
}

func RedisSet(key string, value interface{}, expiration time.Duration) error {
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func RedisGet(key string) (string, error) {
	str, err := rdb.Get(ctx, key).Result()
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	return str, nil
}

func RedisIncr(key string) (int, error) {
	val, err := rdb.IncrBy(ctx, key, 1).Result()
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func RedisDel(key string) {
	rdb.Del(ctx, key)
}
