package middlewares

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type ResponDriver struct {
	Lon       float64 `json:"Lon"`
	Lat       float64 `json:"Lat"`
	Dist      float64 `json:"dist"`
	Driver_id string  `json:"driver_id"`
}

var RedisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6785",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var ctx = context.Background()

func FindDriver(c *fiber.Ctx) error {
	lon := c.QueryFloat("lon")
	lat := c.QueryFloat("lat")
	geo := c.Query("g")[0:4]
	min_km := c.QueryFloat("min_km", 1)

	if len(geo) != 4 || lon == 0 || lat == 0 {
		log.Println(geo, lon, lat)
		return c.SendStatus(400)
	}

	res, err := RedisClient.GeoRadius(ctx, geo, lon, lat, &redis.GeoRadiusQuery{
		Radius: min_km,
		Unit:   "km",
		Count:  10,
		Sort:   "ASC",
	}).Result()
	if err != nil {
		log.Println("Query driver error: ", err)
		return c.SendStatus(500)
	}

	drivers := make([]ResponDriver, 0)
	for _, d := range res {
		drivers = append(drivers, ResponDriver{
			Lon:       d.Longitude,
			Lat:       d.Latitude,
			Dist:      d.Dist,
			Driver_id: d.Name,
		})
	}

	c.SendStatus(200)
	return c.JSON(drivers)
}
