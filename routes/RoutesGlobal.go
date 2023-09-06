package routes

import (
	"context"
	"encoding/json"
	"goserver/libs"
	"io/ioutil"
	"log"
	"net/http"
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

type RideReqInfo struct {
	SLon float64 `json:"slon"`
	SLat float64 `json:"slat"`
	SAdr string  `json:"sadr"`

	ELon float64 `json:"elon"`
	ELat float64 `json:"elat"`
	EAdr string  `json:"eadr"`

	User_id   string `json:"user_id"`
	Driver_id string `json:"driver_id"`
	Trip_id   string `json:"trip_id"`

	Price float64 `json:"price"`

	User_Name string `json:"user_name"`
	User_Phone string `json:"user_phone"`
}
type RideReqToPub struct {
	RideReqInfo
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
	Lon       float64
	Lat       float64
	GeoKey    string
	Driver_id string
}
type DriverLocSSE struct {
	Lon      float64 `json:"lon"`
	Lat      float64 `json:"lat"`
	DriverId string  `json:"driver_id"`
}

var GlobalDriverLocAddChannel = make(chan *DriverLocToAdd, 50)

func RedisAddDriverLocListener() {
	libs.NewPubSub("DriverLoc")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6785",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	context := context.Background()

	for {
		data, ok := <-GlobalDriverLocAddChannel
		if !ok {
			break
		}

		go func(data *DriverLocToAdd) {
			client.GeoAdd(context, data.GeoKey, &redis.GeoLocation{
				Name:      data.Driver_id,
				Longitude: data.Lon,
				Latitude:  data.Lat,
			})
			b, err := json.Marshal(&DriverLocSSE{
				Lon:      data.Lon,
				Lat:      data.Lat,
				DriverId: data.Driver_id,
			})

			if err == nil {
				libs.Publish("DriverLoc", string(b))
			}
		}(data)

	}
}

type ClientDetail struct{
	Phone string `json:"phone"`
	Name string `json:"name"`
}

func GetUserDetailInfo(user_id string) ( *ClientDetail,error) {
	res,err := http.Get("http://localhost:3000/api/users/"+user_id)
	if err != nil{
		return nil,err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil{
		return nil,err
	}

	client_detail := &ClientDetail{
		Phone: "",
		Name: "",
	}
	err = json.Unmarshal(resBody,client_detail)
	if err != nil{
		return nil,err
	}

	return client_detail,nil
}