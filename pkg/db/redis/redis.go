package redis

import (
	"context"
	"time"

	"github.com/gomajido/hospital-cms-golang/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func InitRedis(redisCfg *config.RedisConfig) *Redis {
	if redisCfg.Password == "empty" {
		redisCfg.Password = ""
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Address,  // Update with your Redis server address
		Password: redisCfg.Password, // Set password if required
		DB:       redisCfg.DB,       // Use default database
	})
	return &Redis{
		Client: rdb,
	}
}

func (r *Redis) SetValue(key string, value string, ttl *time.Duration) error {
	var duration time.Duration = 0
	if ttl != nil {
		duration = *ttl
	}

	err := r.Client.Set(context.Background(), key, value, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetValue(key string) (string, error) {
	val, err := r.Client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (r *Redis) GetTTL(key string) (int, error) {
	val, err := r.Client.TTL(context.Background(), key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	durationAsInt := int(val.Seconds()) // Convert to seconds

	return durationAsInt, nil
}
