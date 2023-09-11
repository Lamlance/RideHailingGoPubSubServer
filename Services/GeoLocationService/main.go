package main

import (
	"GoGeoService/middlewares"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	middlewares.ConnectToRedis()
	middlewares.ConnectToRabbitMQ()

	go middlewares.PublishDriverLocation()
	nginx_prefix := "/ridehail/geo"

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Cache-Control,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	
	app.Get(nginx_prefix,func(c *fiber.Ctx) error {return c.SendString("Hello from Geo service")})

	app.Post(nginx_prefix+"/loc/driver/:driver_id", middlewares.DriverLocationPost)
	// app.Get("loc/driver/:driver_id")
	// app.Delete("loc/driver/:driver_id")

	app.Get(nginx_prefix+"/find/drivers", middlewares.FindDriver)

	log.Fatal(app.Listen(":3083"))
}
