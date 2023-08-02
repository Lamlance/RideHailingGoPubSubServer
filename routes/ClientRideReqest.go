package routes

import (
	"goserver/routes/ws"
	"strconv"

	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func publish_ride_request_loop(geo_key string, user_id string, stop_req *chan bool, res []redis.GeoLocation) {
	for i := 0; i < 10; i++ {
		res = append(res, redis.GeoLocation{
			Name: "Driver #" + strconv.Itoa(i),
		})
	}
	stop_loop := make(chan bool, 0)
	timeout_chan := make(chan bool, 0)

	var timer *time.Timer

	timeout_routine := func() {
		//timeout_chan <- false
		log.Println("Start timeout routine")
		if timer == nil {
			timer = time.NewTimer(20 * time.Second)
		} else {
			timer.Reset(20 * time.Second)
		}
		_, ok := <-timer.C
		log.Println("End timeout routine")
		if ok {
			timeout_chan <- true
		}
	}

	go func() {
		driver_ok := false
		for _, pos := range res {
			if driver_ok {
				break
			}
			select {
			case <-(stop_loop):
				log.Println("A driver has stop req loop")
				driver_ok = true
				if !timer.Stop(){
					<- timer.C
				}
				break
			case <-(timeout_chan):
				log.Println("Msg for dirver id: " + pos.Name)
				go timeout_routine()
				GlobalRedisClient.Publish(RedisContext, geo_key, "Msg for dirver id: "+pos.Name)
			default:
				log.Println("Nothing happend")
				time.Sleep(5 * time.Second)
			}
		}
		if !driver_ok {
			*stop_req <- true
			<-stop_loop
		}
		log.Println("Request loop stopped #2")
	}()

	timeout_chan <- true

	_, ok := <-*stop_req
	stop_loop <- true
	if !ok {
		log.Println("Req loop timed out")
	}
	log.Println("Request loop stopped #1")
}

func ClientCheckMiddleware(c *fiber.Ctx) error {
	random_string := "Secrete_trip_id" //uuid.New().String()

	ws.GlobalRoomMap.Lock.Lock()

	room := ws.MakeEmptyCommunicationRoom()
	ws.GlobalRoomMap.Data[random_string] = room

	c.Locals("trip_id", random_string)
	c.Locals("room", room)

	ws.GlobalRoomMap.Lock.Unlock()

	return c.Next()
}

func ClientRideRequest(c *fiber.Ctx) error {
	lon := c.QueryFloat("lon")
	lat := c.QueryFloat("lat")
	user_id := c.Query("user_id")
	geo_hash := c.Params("geo_hash")
	room, ok_room := c.Locals("room").(*ws.CommunicationRoom)

	if len(geo_hash) < 4 {
		c.SendStatus(400)
		return c.SendString("Invalid geo hash")
	}

	if !ok_room {
		log.Println("Server can't find communication room")
		return c.SendStatus(500)
	}

	geo_key := geo_hash[0:4]

	res, err := GlobalRedisClient.GeoRadius(RedisContext, geo_key, lon, lat, &redis.GeoRadiusQuery{
		Radius: 1,
		Unit:   "km",
		Count:  10,
		Sort:   "ASC",
	}).Result()

	if err != nil {
		log.Println("Ride request error: ", err)
		return c.SendStatus(500)
	}
	if len(res) <= 0 {
		c.SendString("Server can't find any driver")
		return c.SendStatus(500)
	}
	go publish_ride_request_loop(geo_key, user_id, &room.Break_ride_request, res)
	return c.Next()
}
