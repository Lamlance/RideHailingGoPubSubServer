package main

import (
	"goserver/libs"
	"goserver/routes/ws"
	"goserver/routes/xhr"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	app.Get("/ws/client",
		ws.ClientCheckMiddleware,
		websocket.New(ws.ClientListenThread))

	app.Get("/ws/driver/:trip_id",
		ws.DriverHandlerMiddleware,
		websocket.New(ws.DriverListenThread))

	app.Get("xhr/driver/",
		xhr.DriverWaitRequest)

	go libs.KafkaConsumer()

	log.Fatal(app.Listen(":3080"))
}
