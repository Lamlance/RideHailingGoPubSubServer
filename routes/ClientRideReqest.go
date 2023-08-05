package routes

import (
	"goserver/routes/ws"
	"strconv"

	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func publish_ride_request_loop(geo_key string, user_id string, req_chan *chan int, res []redis.GeoLocation) {
	for i := 0; i < 3; i++ {
		res = append(res, redis.GeoLocation{
			Name: "Driver #" + strconv.Itoa(i),
		})
	}

	req_code := make(chan int, 0)
	timeout_chan := make(chan bool, 0)
	skip_timeout_chan := make(chan bool, 0)

	var timer *time.Timer
	timeout_routine := func() {
		log.Println("Start timeout routine")
		if timer == nil {
			timer = time.NewTimer(20 * time.Second)
		} else {
			timer.Reset(20 * time.Second)
		}
		select {
		case <-skip_timeout_chan:
			log.Println("Time out has been skipped")

			timeout_chan <- true
		case <-timer.C:
			log.Println("End timeout routine")
			timeout_chan <- true
		}

	}

	go func() {
		driver_ok := false

		for i := 0; i < len(res)+1; i++ {
			if driver_ok || i >= len(res) {
				break
			}

			pos := res[i]
			select {
			case res := <-(req_code):
				if !timer.Stop() {
					log.Println("Stop timeout false")
					<-timer.C
				} else {
					log.Println("Stop timeout true")
				}

				if res == 0 {
					log.Println("A driver has stop req loop")
					driver_ok = true
				} else {
					skip_timeout_chan <- true
					log.Println("A Driver has skipped request")
				}

			case <-(timeout_chan):
				log.Println("Msg for dirver id: " + pos.Name)
				if i < len(res) -1 {
					go timeout_routine()
				}
				GlobalRedisClient.Publish(RedisContext, geo_key, "Msg for dirver id: "+pos.Name)
			}
		}
		log.Println("Request loop stopping #2")

		if !driver_ok {
			log.Println("No driver accept request")
			*req_chan <- 3
			<-req_code
		}

		log.Println("Request loop stopped #2")
	}()

	timeout_chan <- true
	for {
		code, ok := <-*req_chan
		if !ok {
			log.Println("Req loop timed out")
			break
		}
		req_code <- code
		if code == 0 || code == 3 {
			break
		}
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
	go publish_ride_request_loop(geo_key, user_id, &room.Ride_requst_channel, res)
	return c.Next()
}
