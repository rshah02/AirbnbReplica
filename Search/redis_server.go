package main

import (
  "github.com/go-redis/redis"
)

func NewRedisServer() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis.sc2jno.ng.0001.usw2.cache.amazonaws.com:6379",
		Password: "",
		DB:       0,  // use default DB
	})

  return client
