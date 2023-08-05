package routes

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var GlobalRedisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6785",
	Password: "", // no password set
	DB:       0,  // use default DB
})
var RedisContext = context.Background()

func RedisSubscribe(topics []string) {

	sub := GlobalRedisClient.Subscribe(RedisContext, topics...)
	defer sub.Close()

	for {
		msg, ok := <-sub.Channel()
		if !ok {
			log.Println("Error reading redis subscribe")
			break
		} else {
			log.Println("Message: '" + msg.Payload + "' from channel: " + msg.Channel)
		}
	}
}

func FindNearestDriver(lon float64, lat float64, geo_key string) ([]redis.GeoLocation, error) {
	res, err := GlobalRedisClient.GeoRadius(RedisContext, geo_key, lon, lat, &redis.GeoRadiusQuery{
		Radius: 1,
		Unit:   "km",
		Count:  10,
		Sort:   "ASC",
	}).Result()
	return res, err
}

func ride_req_loop(drivers []redis.GeoLocation) {
	skip_chan := make(chan bool)
	accept_chan := make(chan bool)

	var timer *time.Timer

	timeout_routine := func() {
		if timer == nil {
			timer = time.NewTimer(20 * time.Second)
		} else {
			timer.Reset(20 * time.Second)
		}
		<-timer.C
		skip_chan <- true
	}

	go func() {
		driver_ok := false

		for _, d := range drivers {
			if driver_ok {
				break
			}

			select {
			case <-accept_chan:
				log.Println("A driver had accepted ride")
				if !timer.Stop() {
					<-timer.C
				}
			case <-skip_chan:
				log.Println("Publish message for driver: ", d.Name)
				if !timer.Stop() {
					<-timer.C
				}
				go timeout_routine()
			}
		}
	}()
}
