package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

var redisClient *RedisClient

type RedisClient struct {
	prefix string

	rdb    *redis.Client
	locker *redislock.Client
}

func GetRedisClient() *RedisClient {
	return redisClient
}

func NewRedisClient(rdb *redis.Client, prefix string) *RedisClient {
	redisClient = &RedisClient{rdb: rdb, prefix: prefix, locker: redislock.New(rdb)}
	return redisClient
}

// Set save a (key,value) to redis
func (repo *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return repo.rdb.Set(context.Background(), repo.keyPrefix(key), value, expiration).Err()
}

// GetString return a string value from redis by specisal key
func (repo *RedisClient) GetString(key string) (string, error) {
	return repo.rdb.Get(context.Background(), repo.keyPrefix(key)).Result()
}

// GetInt64 return a int64 value from redis by specisal key
func (repo *RedisClient) GetInt64(key string) (int64, error) {
	return repo.rdb.Get(context.Background(), repo.keyPrefix(key)).Int64()
}

// SetObject save a object value to redis by specisal key
func (repo *RedisClient) SetObject(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return repo.rdb.Set(context.Background(), repo.keyPrefix(key), data, expiration).Err()
}

// GetObject return a object value from redis by specisal key
func (repo *RedisClient) GetObject(key string, value interface{}) error {
	data, err := repo.rdb.Get(context.Background(), repo.keyPrefix(key)).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

// HSet save a simple value to redis hash table by specisal key
func (repo *RedisClient) HSet(key string, values ...interface{}) error {
	return repo.rdb.HSet(context.Background(), repo.keyPrefix(key), values...).Err()
}

// HExists returns whether the field exists
func (repo *RedisClient) HExists(key string, field string) (bool, error) {
	if key == "" || field == "" {
		return false, errors.New("参数为空")
	}
	has, err := repo.rdb.HExists(context.Background(), repo.keyPrefix(key), field).Result()
	if err != nil {
		return false, errors.New("异常：" + err.Error())
	}
	return has, nil
}

// Delete remove a value from redis by specisal key
func (repo *RedisClient) Delete(key string) error {
	return repo.rdb.Del(context.Background(), repo.keyPrefix(key)).Err()
}

// Obtain try to obtain lock.
func (repo *RedisClient) Obtain(ctx context.Context, key string, ttl time.Duration) (*redislock.Lock, error) {
	return repo.locker.Obtain(ctx, repo.keyPrefix(key), ttl, nil)
}

func (repo *RedisClient) keyPrefix(key string) string {
	return fmt.Sprintf("%s:%s", repo.prefix, key)
}

func (repo *RedisClient) Incr(key string) int64 {
	return repo.rdb.Incr(context.Background(), repo.keyPrefix(key)).Val()
}

func (repo *RedisClient) ZAdd(ctx context.Context, key string, members ...*redis.Z) error {
	return repo.rdb.ZAdd(ctx, repo.keyPrefix(key), members...).Err()
}

func (repo *RedisClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	r, err := repo.rdb.ZRange(ctx, repo.keyPrefix(key), start, stop).Result()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (repo *RedisClient) ZCard(ctx context.Context, key string) (int64, error) {
	r, err := repo.rdb.ZCard(ctx, repo.keyPrefix(key)).Result()
	if err != nil {
		return 0, err
	}

	return r, nil
}
