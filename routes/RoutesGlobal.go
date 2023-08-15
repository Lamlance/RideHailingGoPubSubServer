package routes

import (
	"context"
	"encoding/json"
	"goserver/libs"
	"goserver/routes/ws"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var GlobalRedis struct {
	Client  *redis.Client
	Context context.Context
	Mutex   *sync.Mutex
} = struct {
	Client  *redis.Client
	Context context.Context
	Mutex   *sync.Mutex
}{
	Client: redis.NewClient(&redis.Options{
		Addr:     "localhost:6785",
		Password: "", // no password set
		DB:       0,  // use default DB
	}),
	Context: context.Background(),
	Mutex:   new(sync.Mutex),
}

func RedisSubscribe(topics []string) {

	GlobalRedis.Mutex.Lock()
	sub := GlobalRedis.Client.Subscribe(GlobalRedis.Context, topics...)
	defer sub.Close()
	GlobalRedis.Mutex.Unlock()

	for {
		msg, ok := <-sub.Channel()
		if !ok {
			log.Println("Error reading redis subscribe")
			break
		} else {
			//log.Println("Message: '" + msg.Payload + "' from channel: " + msg.Channel)
			libs.Publish(msg.Channel, msg.Payload)
		}
	}
}

type RideReqToPub struct {
	ws.RideReqInfo
	Channel string
}

var GlobalRideReqToPubChannel = make(chan *RideReqToPub, 50)

func RedisPublishRideReqListener() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6785",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	context := context.Background()

	for {
		data, ok := <-GlobalRideReqToPubChannel
		if !ok {
			break
		}
		go func(req *RideReqToPub) {
			b, _ := json.Marshal(data)
			client.Publish(context, req.Channel, b)
		}(data)

	}
}

type DriverLocToAdd struct {
	Lon float64
	Lat float64
	GeoKey string
	Driver_id string
}

var GlobalDriverLocAddChannel = make(chan *DriverLocToAdd, 50)

func RedisAddDriverLocListener() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6785",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	context := context.Background()

	for{
		data,ok := <- GlobalDriverLocAddChannel
		if !ok {
			break
		}

		go func (data *DriverLocToAdd)  {
			client.GeoAdd(context,data.GeoKey,&redis.GeoLocation{
				Name: data.Driver_id,
				Longitude: data.Lon,
				Latitude: data.Lat,
			})
		}(data)

	}
}
