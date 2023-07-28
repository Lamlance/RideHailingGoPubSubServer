package loc

import (
	"goserver/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type DriverLocationPostBody struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
	G   string  `json:"g"`
}

func DriverLocationPost(c *fiber.Ctx) error {

	driver_id := c.Params("driver_id")
	if driver_id == "" {
		return c.SendStatus(404)
	}

	body := DriverLocationPostBody{
		Lon: 0,
		Lat: 0,
	}

	if err := c.BodyParser(&body); err != nil {
		log.Println("Body praser error: ", err)
		return c.SendStatus(400)
	} else {
		log.Println("Driver location body: ", body)
	}

	_,err := routes.GlobalRedisClient.GeoAdd(routes.RedisContext, body.G, &redis.GeoLocation{
		Longitude: body.Lon,
		Latitude:  body.Lat,
		Name:      driver_id,
	}).Result()

	if err != nil{
		log.Println("Geo add error: ",err)
		return c.SendStatus(500)
	}

	return c.SendStatus(202)
}

func DriverLocationGet(c *fiber.Ctx) error {
	driver_id := c.Params("driver_id")
	geo_hash := c.Query("geo_hash")
	if driver_id == "" || geo_hash == ""{
		return c.SendStatus(404)
	}

	res,err := routes.GlobalRedisClient.GeoPos(routes.RedisContext,geo_hash,driver_id).Result()
	if err != nil{
		log.Println("Get geo area: ",err)
		return c.SendStatus(500)
	}
	
	for _,pos := range res{
		log.Println("Location: ",pos.Longitude,pos.Latitude)
	}

	return c.SendStatus(200)
}
