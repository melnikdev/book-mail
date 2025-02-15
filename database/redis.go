package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		fmt.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}

	return db, nil
}
