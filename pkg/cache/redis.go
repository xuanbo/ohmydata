package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/xuanbo/ohmydata/pkg/config"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
)

// Init 初始化
func Init() error {
	addr := config.GetString("redis.addr")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	db := config.GetInt("redis.db")
	password := config.GetString("redis.password")

	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Password: password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return err
	}

	return nil
}

// Set 设置缓存，小于1s默认为1h
func Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if ttl < time.Second {
		ttl = time.Hour
	}
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return client.Set(ctx, key, b, ttl).Err()
}

// Get 查询缓存
func Get(ctx context.Context, key string, value interface{}) error {
	b, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, value)
}

// Del 清除缓存
func Del(ctx context.Context, keys ...string) error {
	return client.Del(ctx, keys...).Err()
}

// DelMatch 清除缓存
func DelMatch(ctx context.Context, match string) error {
	var (
		cursor uint64
		count  int64 = 100
	)
	for {
		var (
			keys []string
			err  error
		)
		keys, cursor, err = client.Scan(ctx, cursor, match, count).Result()
		if err != nil {
			return err
		}
		// 清除
		if err := Del(ctx, keys...); err != nil {
			return err
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}
