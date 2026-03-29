package cache

import (
	"context"
	"fmt"

	"github.com/2yuri/pack-calculator/pkg/config"
	"github.com/redis/go-redis/v9"
)

func New() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Instance().Cache.Host, config.Instance().Cache.Port),
		Password: config.Instance().Cache.Password,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
