package main

import (
	"GoTripService/middlewares"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	go middlewares.PublishRideRequest()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Cache-Control,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Get("/ridehail/trip/", func(c *fiber.Ctx) error { return c.SendString("Hello Trip service") })

	//Client side
	app.Get("/ridehail/trip/ws/client/:geo_hash",
		middlewares.TripMiddleware,
		middlewares.ClientRideRequest,
		websocket.New(middlewares.ClientListenThread),
	)

	//Driver side
	app.Get("/ridehail/trip/ws/driver/:trip_id",
		middlewares.DriverRideRequest,
		websocket.New(middlewares.DriverListenThread),
	)

	//Admin side
	app.Get("/ridehail/trip/admin/client/:geo_hash",
		middlewares.TripMiddleware,
		middlewares.ClientRideRequest,
		websocket.New(middlewares.AdminRoomHandler))

	log.Fatal(app.Listen(":3081"))

}
