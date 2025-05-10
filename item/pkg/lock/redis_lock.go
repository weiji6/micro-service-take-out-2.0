package lock

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLock struct {
	client     *redis.Client
	key        string
	value      string
	expiration time.Duration
}

func NewRedisLock(client *redis.Client, key string, expiration time.Duration) *RedisLock {
	return &RedisLock{
		client:     client,
		key:        key,
		value:      time.Now().String(), // 使用时间戳作为锁的值
		expiration: expiration,
	}
}

// TryLock 尝试获取锁
func (l *RedisLock) TryLock(ctx context.Context) (bool, error) {
	return l.client.SetNX(ctx, l.key, l.value, l.expiration).Result()
}

// Unlock 释放锁
func (l *RedisLock) Unlock(ctx context.Context) error {
	// 使用Lua脚本确保原子性操作
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	return l.client.Eval(ctx, script, []string{l.key}, l.value).Err()
}

// AutoLock 自动获取锁并在指定时间后释放
func (l *RedisLock) AutoLock(ctx context.Context, fn func() error) error {
	locked, err := l.TryLock(ctx)
	if err != nil {
		return err
	}
	if !locked {
		return ErrLockFailed
	}

	defer l.Unlock(ctx)
	return fn()
} 