package routes

import (
	"context"
	"goserver/libs"
	"log"

	"github.com/redis/go-redis/v9"
)

var GlobalRedisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6785",
	Password: "", // no password set
	DB:       0,  // use default DB
})
var RedisContext = context.Background()

func RedisSubscribe() {
	sub := GlobalRedisClient.Subscribe(RedisContext, "w3gv")
	defer sub.Close()

	for {
		msg, ok := <-sub.Channel()
		if !ok {
			log.Println("Error reading redis subscribe")
		} else {
			log.Println("Message: '" + msg.Payload + "' from channel: " + msg.Channel)
			libs.Publish(msg.Channel,msg.Payload)
		}
	}

}
