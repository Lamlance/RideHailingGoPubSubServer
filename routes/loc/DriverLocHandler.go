package loc

import (
	"goserver/routes"
	"log"
	"strconv"

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

	if len(body.G) < 4{
		c.SendStatus(400)
		return c.SendString("Invalid geo hash")
	}

	geo_key := body.G[0:3]

	_, err := routes.GlobalRedisClient.GeoAdd(routes.RedisContext, geo_key , &redis.GeoLocation{
		Longitude: body.Lon,
		Latitude:  body.Lat,
		Name:      driver_id,
	}).Result()

	if err != nil {
		log.Println("Geo add error: ", err)
		return c.SendStatus(500)
	}

	return c.SendStatus(202)
}

func DriverLocationGet(c *fiber.Ctx) error {
	driver_id := c.Params("driver_id")
	geo_hash := c.Query("geo_hash")
	if driver_id == "" || len(geo_hash) < 4 {
		return c.SendStatus(404)
	}
	
	geo_key := geo_hash[0:3]

	res, err := routes.GlobalRedisClient.GeoSearchLocation(routes.RedisContext,geo_key,
		&redis.GeoSearchLocationQuery{
			GeoSearchQuery: redis.GeoSearchQuery{
				Count: 1,
				Member: driver_id,
			},
			WithCoord: true,
		}).Result()

	if err != nil {
		log.Println("Get geo error: ", err)
		return c.SendStatus(500)
	}

	if len(res) <= 0{
		c.SendStatus(404)
		return c.SendString("Driver not found")
	}

	resJson := struct {
		Lon float64 `json:"Lon"`
		Lat float64 `json:"Lat"`
	}{
		Lon: res[0].Longitude,
		Lat: res[0].Latitude,
	}

	c.SendStatus(200)
	return c.JSON(resJson)
}

func DriverLocationDelete(c *fiber.Ctx) error {
	driver_id := c.Params("driver_id")
	geo_hash := c.Query("geo_hash")
	if driver_id == "" || len(geo_hash) < 4 {
		return c.SendStatus(404)
	}

	geo_key := geo_hash[0:3]

	res,err := routes.GlobalRedisClient.ZRem(routes.RedisContext, geo_key, driver_id).Result()

	if err != nil{
		c.SendStatus(500)
		return c.SendString(err.Error())
	}

	c.SendStatus(200)
	return c.SendString(strconv.FormatInt(res,10))
}
