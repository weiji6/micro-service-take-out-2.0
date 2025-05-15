package lock

import (
	"context"

	"github.com/go-redsync/redsync/v4"
	redsyncgoredis "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	Rdb     *redis.Client
	RedSync *redsync.Redsync
)

func InitRedisLock() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	if err := Rdb.Ping(context.Background()).Err(); err != nil {
		panic("无法连接 Redis:" + err.Error())
	}

	pool := redsyncgoredis.NewPool(Rdb)
	RedSync = redsync.New(pool)
}
