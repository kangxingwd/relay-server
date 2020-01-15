package reids

import (
	"github.com/go-redis/redis"
	cfg "gitlab.com/TenbayMCloud/awesome-raserver/config"
)

var redisClient *redis.Client

func RedisConnectInit() error {
	c := redis.NewClient(&redis.Options{
		Addr:		cfg.GetRedis().Address,
		Password:	cfg.GetRedis().Password,
		DB:			cfg.GetRedis().AccountLockDB,
	})

	_, err := c.Ping().Result()
	if err != nil {
		return err
	}
	redisClient = c
	return nil
}
func RedisClose() error {
	return redisClient.Close()
}
func RedisClient() *redis.Client {
	return redisClient
}
