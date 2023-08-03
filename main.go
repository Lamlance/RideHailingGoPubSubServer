package main

import (
	"goserver/libs"
	"goserver/routes"
	"goserver/routes/loc"
	"goserver/routes/ws"
	"goserver/routes/xhr"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Get("/ws/client/:geo_hash",
		routes.ClientCheckMiddleware,
		routes.ClientRideRequest,
		websocket.New(ws.ClientListenThread))

	app.Get("/ws/driver/:trip_id",
		ws.DriverHandlerMiddleware,
		websocket.New(ws.DriverListenThread))

	app.Get("xhr/driver/:geo_hash",
		xhr.DriverWaitRequest)

	app.Post("loc/driver/:driver_id", loc.DriverLocationPost)
	app.Get("loc/driver/:driver_id", loc.DriverLocationGet)
	app.Delete("loc/driver/:driver_id", loc.DriverLocationDelete)

	libs.NewPubSub("w3gv")

	//go libs.KafkaConsumer()
	go routes.RedisSubscribe()

	log.Fatal(app.Listen(":3080"))
}
