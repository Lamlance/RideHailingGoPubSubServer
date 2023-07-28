package routes

import (
	"context"
	"goserver/libs"

	"github.com/redis/go-redis/v9"
)

var PubSub = libs.NewPubSub("Driver")
var GlobalRedisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6785",
	Password: "", // no password set
	DB:       0,  // use default DB
})
var RedisContext = context.Background()
