package initialization

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v8"

	"com.say.more.server/config"
)

// Redis connect tht redis server
func Redis(cfg *config.Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.Db,
	})
}

// RedisClose Has return a bool value
func RedisClose(rdb *redis.Client, prefix string) {
	p := fmt.Sprintf("%s*", prefix)

	ctx := context.Background()
	iter := rdb.Scan(ctx, 0, p, 0).Iterator()
	for iter.Next(ctx) {
		err := rdb.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	_ = rdb.Close()
}
