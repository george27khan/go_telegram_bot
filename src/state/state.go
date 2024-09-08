package state

import (
	"context"
	"go_telegram_bot/src/redis"
)

func Set(ctx context.Context, state_key string, idUser string, value string) {
	rdb := redis.Connect()
	defer rdb.Close()
	_ = rdb.HSet(ctx, "user_state", map[string]interface{}{idUser: value})
}

func Get(ctx context.Context, state_key string, idUser string) string {
	rdb := redis.Connect()
	defer rdb.Close()
	res := rdb.HGet(ctx, state_key, idUser)
	return res.Val()
}

func Del(ctx context.Context, state_key string, idUser string) {
	rdb := redis.Connect()
	defer rdb.Close()
	_ = rdb.HDel(ctx, state_key, idUser)
}
