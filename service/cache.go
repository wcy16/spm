// use redis as cache
package service

import (
	"github.com/go-redis/redis"
	"log"
)

const (
	Expire1min  = 60
	Expire10min = 600
	Expire30min = 1800
)

func init() {
	register(connectToRedis)
}

var Cache *redis.Client

func connectToRedis() {
	Cache = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// test connection
	if _, err := Cache.Ping().Result(); err != nil {
		log.Println("Redis not initialized.")
	} else {
		log.Println("Redis connected.")
	}
}
