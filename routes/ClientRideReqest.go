package routes

import (
	"goserver/routes/ws"

	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func publish_ride_request_loop(geo_key string, user_id string, stop_req *chan bool, res []redis.GeoLocation) {
	stop_loop := make(chan bool, 0)

	go func() {
		for _, pos := range res {
			select {
			case <-(stop_loop):
				log.Println("A driver has stop req loop")
				break
			default:
				log.Println("Send ride request for driver")
				
				GlobalRedisClient.Publish(RedisContext,geo_key,"Msg for dirver id: " + pos.Name)

				time.Sleep(2 * time.Second)
			}
		}
		<-stop_loop
		log.Println("Request loop stopped #2")

	}()
	_, ok := <-*stop_req
	stop_loop <- true
	if !ok {
		log.Println("Req loop timed out")
	}
	log.Println("Request loop stopped #1")

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
	if len(res) <= 0{
		c.SendString("Server can't find any driver")
		return c.SendStatus(500)
	}
	go publish_ride_request_loop(geo_key, user_id, &room.Break_ride_request, res)
	return c.Next()
}
