package main

import (
	"GoSSEService/middlewares"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	middlewares.NewPubSub("w3gv")
	middlewares.NewPubSub("DriverCoord")

	go middlewares.ListenRideRequest()
	go middlewares.ListenDriverLoc()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Cache-Control,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("Hello SSE service") })

	app.Get("xhr/driver/:geo_hash", middlewares.DriverListenRideRequest)
	app.Get("sse/driver_loc/:driver_id", middlewares.DriverLocationWatch)

	log.Fatal(app.Listen(":3082"))
}
