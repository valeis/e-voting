package web

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

// ChainCodeQuery interface for querying chaincode -> Component Interface
type ChainCodeQuery interface {
	Query(chainCodeName, channelID, function string, args []string) (string, error)
}

// SimpleQuery implements the base chaincode query functionality -> Concrete Components
type SimpleQuery struct {
	setup OrgSetup
}

func (sq *SimpleQuery) Query(chainCodeName, channelID, function string, args []string) (string, error) {
	network := sq.setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.EvaluateTransaction(function, args...)
	if err != nil {
		return "", err
	}
	return string(evaluateResponse), nil
}

// Decorator Structs
type RedisCacheQueryDecorator struct {
	query       ChainCodeQuery
	redisClient *redis.Client
	ctx         context.Context
	ttl         time.Duration
}

func NewRedisCacheQueryDecorator(query ChainCodeQuery, redisClient *redis.Client, ttl time.Duration) *RedisCacheQueryDecorator {
	ctx := context.Background()

	return &RedisCacheQueryDecorator{
		query:       query,
		redisClient: redisClient,
		ctx:         ctx,
		ttl:         ttl,
	}
}

// RedisCacheQueryDecorator
func (rcq *RedisCacheQueryDecorator) Query(chainCodeName, channelID, function string, args []string) (string, error) {
	key := fmt.Sprintf("query:%s:%s:%s:%v", channelID, chainCodeName, function, args)

	cachedResult, err := rcq.redisClient.Get(rcq.ctx, key).Result()
	if err == nil {
		fmt.Printf("Cache hit for key: %s", key)
		return cachedResult, nil
	} else if err != redis.Nil {
		log.Printf("Redis error: %v", err)
	}

	result, err := rcq.query.Query(chainCodeName, channelID, function, args)
	if err != nil {
		return "", err
	}
	err = rcq.redisClient.Set(rcq.ctx, key, result, rcq.ttl).Err()
	if err != nil {
		log.Printf("Failed to cache result in Redis: %v", err)
	}
	return result, nil
}

func (setup OrgSetup) QueryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")
	function := queryParams.Get("function")
	args := r.URL.Query()["args"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)

	baseQuery := &SimpleQuery{setup: setup}

	cacheTTL := 5 * time.Minute
	redisClient, ok := setup.RedisClient.(*redis.Client)
	if !ok {
		return
	}
	cachedQuery := NewRedisCacheQueryDecorator(baseQuery, redisClient, cacheTTL)

	response, err := cachedQuery.Query(chainCodeName, channelID, function, args)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Response: %s", response)
}
