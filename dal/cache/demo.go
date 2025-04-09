package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dog-frame/common/enum"
	"github.com/dog-frame/common/logger"
	"github.com/dog-frame/logic/do"
)

// 关于使用 HSET 直接存储结构体的讨论 https://github.com/redis/go-redis/discussions/2454

type DummyDemo struct {
	Name string `redis:"name"`
	ID   int64  `redis:"id"`
}

// SetDemoStruct 使用HSET的存储结构体数据
func SetDemoStruct(ctx context.Context, Demo *do.Demo) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO, Demo.Name)
	data := struct {
		Name string `redis:"Name"`
		ID   int64  `redis:"ID"`
	}{
		ID:   Demo.ID,
		Name: Demo.Name,
	}
	_, err := Redis().HSet(ctx, redisKey, data).Result()
	if err != nil {
		logger.Error(ctx, "redis error", "err", err)
		return err
	}

	return nil
}

// GetDemoStruct 使用HGETALL 和 Scan 读取结构体数据
func GetDemoStruct(ctx context.Context, Name string) (*DummyDemo, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO, Name)
	data := new(DummyDemo)
	err := Redis().HGetAll(ctx, redisKey).Scan(&data)
	Redis().Get(ctx, redisKey).String()
	if err != nil {
		logger.Error(ctx, "redis error", "err", err)
		return nil, err
	}
	logger.Info(ctx, "scan data from redis", "data", &data)
	return data, nil

}

func SetDemo(ctx context.Context, Demo *do.Demo) error {
	jsonDataBytes, _ := json.Marshal(Demo)
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO, Demo.Name)
	_, err := Redis().Set(ctx, redisKey, jsonDataBytes, 0).Result()
	if err != nil {
		logger.Error(ctx, "redis error", "err", err)
		return err
	}

	return nil
}

func GetDemo(ctx context.Context, Name string) (*do.Demo, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO, Name)
	jsonBytes, err := Redis().Get(ctx, redisKey).Bytes()
	if err != nil {
		logger.Error(ctx, "redis error", "err", err)
		return nil, err
	}
	data := new(do.Demo)
	err = json.Unmarshal(jsonBytes, &data)
	return data, err
}
