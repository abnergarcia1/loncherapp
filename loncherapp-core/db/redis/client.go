package redis

import (
	"os"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var Client *redis.Client

func InitRedis() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	Client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func NewRedisDB(dsn, password string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: password,
		DB:       0,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.WithFields(log.Fields{
			"RedisClient": redisClient,
		}).Errorf("Error when trying to get redis db: %v", err)
	}
	return redisClient
}
