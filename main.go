package main

import (
	"goserver/libs"
	"goserver/routes"
	"goserver/routes/loc"
	"goserver/routes/sse"
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
		AllowHeaders:     "Origin,Cache-Control,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Get("/", func(c *fiber.Ctx) error { return c.SendString("Hello") })

	app.Get("/ws/client/:geo_hash",
		ws.ClientCheckMiddleware,
		ws.ClientRideRequest,
		websocket.New(ws.ClientListenThread))

	app.Get("/admin/client/:geo_hash",
		ws.ClientCheckMiddleware,
		ws.ClientRideRequest,
		websocket.New(ws.AdminRoomHandler))

	app.Get("/ws/driver/:trip_id",
		ws.DriverHandlerMiddleware,
		websocket.New(ws.DriverListenThread))

	app.Delete("/ws/driver/:trip_id",
		ws.DriverHandlerMiddleware,
		func(c *fiber.Ctx) error {
			room, ok_room := c.Locals("room").(*ws.CommunicationRoom)

			if !ok_room {
				log.Println("Driver cant find com channel")
				return c.SendStatus(500)
			}

			log.Printf("Driver is skipping ride req loop")
			room.Ride_requst_channel <- 1
			log.Printf("Driver has skipping ride req loop")
			return c.SendStatus(200)
		})

	app.Get("xhr/driver/:geo_hash",
		xhr.DriverWaitRequest)

	app.Post("loc/driver/:driver_id", loc.DriverLocationPost)
	app.Get("loc/driver/:driver_id", loc.DriverLocationGet)
	app.Delete("loc/driver/:driver_id", loc.DriverLocationDelete)

	app.Get("sse/driver_loc/:driver_id", sse.DriverLoc)
	app.Get("sse/driver_wait/:geo_hash", sse.DriverWaitReq)

	libs.NewPubSub("w3gv")
	
	//go libs.KafkaConsumer()
	topics := []string{"w3gv"}
	go routes.RedisSubscribe(topics)
	go routes.RedisPublishRideReqListener()
	go routes.RedisAddDriverLocListener()

	log.Fatal(app.Listen(":3080"))
}
