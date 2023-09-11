package middlewares

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type DriverLocationPostBody struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
	G   string  `json:"g"`
}

func DriverLocationPost(c *fiber.Ctx) error {
	c.Request().Header.SetContentType("application/json")
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

	if len(body.G) < 4 {
		c.SendStatus(400)
		return c.SendString("Invalid geo hash")
	}

	RedisUpdateDriver(body.Lon, body.Lat, body.G, driver_id)

	DriverLocationToPub <- &DriverPubLoc{
		Lon:       body.Lon,
		Lat:       body.Lat,
		Driver_id: driver_id,
	}

	return c.SendStatus(200)
}
