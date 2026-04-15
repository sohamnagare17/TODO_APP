package Redis 

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)


var Ctx = context.Background()

func GetCache(rdb *redis.Client, key string, dest interface{}) bool {
	val, err := rdb.Get(Ctx, key).Result()
	if err != nil {
		return false
	}
	json.Unmarshal([]byte(val), dest)
	return true
}

func SetCache(rdb *redis.Client, key string, data interface{}, ttl time.Duration) {
	bytes, _ := json.Marshal(data)
	rdb.Set(Ctx, key, bytes, ttl)
}

func InitRedis() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis address
	})

	return rdb
}