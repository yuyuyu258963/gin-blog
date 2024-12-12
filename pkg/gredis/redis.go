package gredis

import (
	"context"
	"encoding/json"
	"gin_example/pkg/logging"
	"gin_example/pkg/setting"
	"time"

	"log"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

// Redis 数据库的相关配置
func Setup() {
	RedisSetting := setting.RedisSetting

	rdb = redis.NewClient(&redis.Options{
		Addr:           RedisSetting.Addr,
		Password:       RedisSetting.Password,
		DB:             RedisSetting.DB,
		DialTimeout:    RedisSetting.ConnectTimeout,
		ReadTimeout:    RedisSetting.ReadTimeout,
		WriteTimeout:   RedisSetting.WriteTimeout,
		MinIdleConns:   RedisSetting.MinIdleConns,
		MaxIdleConns:   RedisSetting.MaxIdleConns,
		MaxActiveConns: RedisSetting.MaxActiveConns,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		logging.FatalF("failed to connect redis : %+v", rdb)
		panic(err)
	}

	logging.InfoF("success connect to redis: %s %d", RedisSetting.Addr, RedisSetting.DB)
	log.Printf("success connect to redis: %s %d", RedisSetting.Addr, RedisSetting.DB)
}

// 设置一个key并指定过期时间
func Set(key string, data interface{}, seconds int) error {
	ctx := context.Background()
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// SetEx 直接设置值和过期时间
	return rdb.SetEx(ctx, key, value,
		time.Duration(seconds)*time.Second).Err()
}

// 检查Key对应的值是否存在
func Exists(key string) bool {
	ctx := context.Background()

	exists, err := rdb.Exists(ctx, key).Result()
	// Exists 返回值是int64类型，表示存在的key的数量
	if err != nil {
		return false
	}
	return exists > 0
}

// 获取一个key对应的值
func Get(key string) (string, error) {
	ctx := context.Background() // 可以设置超时时间

	res, err := rdb.Get(ctx, key).Result()

	switch {
	case err == redis.Nil:
		return "", nil
	case err != nil:
		logging.ErrorF("failed redis get key: %v", err)
		return "", err
	default:
		return res, nil
	}
}

// 删除一个key对应的值
func Delete(key string) error {
	ctx := context.Background()
	return rdb.Del(ctx, key).Err()
}

// 模糊删除所有的key
func LikeDeletes(pattern string) error {
	ctx := context.Background()

	// 使用SCAN 代替KEYS
	iter := rdb.Scan(ctx, 0, "*"+pattern+"*", 0).Iterator()

	// 使用Pipeline批量删除
	pipe := rdb.Pipeline()

	for iter.Next(ctx) {
		pipe.Del(ctx, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	// 执行 Pipeline
	_, err := pipe.Exec(ctx)
	return err
}
