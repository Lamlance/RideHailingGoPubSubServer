package main

import (
	"GoGeoService/middlewares"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	go middlewares.PublishDriverLocation()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Cache-Control,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Post("/ridehail/geo/loc/driver/:driver_id", middlewares.DriverLocationPost)
	// app.Get("loc/driver/:driver_id")
	// app.Delete("loc/driver/:driver_id")

	app.Get("/ridehail/geo/find/drivers", middlewares.FindDriver)

	log.Fatal(app.Listen(":3083"))
}
