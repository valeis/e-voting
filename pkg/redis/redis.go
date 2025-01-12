package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

type ConnectionRedis struct {
	ConnectionString string
}

func (con ConnectionRedis) Connect() any {
	opt, err := redis.ParseURL(con.ConnectionString)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)

	if client == nil {
		fmt.Println("failed to connect to RedisDB")
	}
	fmt.Println("successfully connected to RedisDB")
	return client
}
