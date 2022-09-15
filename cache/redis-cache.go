package cache

import (
	"encoding/json"
	"time"

	redis "github.com/go-redis/redis/v7"
)

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) PostCache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: exp,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, value string) {
	client := cache.getClient()
	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	client.Set(key, json, cache.expires*time.Second)
}

func (cache *redisCache) Get(key string) string {
	client := cache.getClient()

	keyValue, err := client.Get(key).Result()
	if err != nil {
		return ""
	}
	var value string
	err = json.Unmarshal([]byte(keyValue), &value)
	if err != nil {
		panic(err)
	}
	return value

}
