package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func RedisPooling(addr, password string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, //default db

		//pool ko settings
		PoolSize:        64,
		MaxIdleConns:    4,
		PoolTimeout:     30 * time.Second,
		ConnMaxIdleTime: 240 * time.Second,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}
	return rdb
}
