package main

import (
	"goserver/routes"
	"goserver/routes/loc"
	"goserver/routes/ws"
	"goserver/routes/xhr"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

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

	//go libs.KafkaConsumer()
	go routes.RedisSubscribe()
	log.Fatal(app.Listen(":3080"))
}
