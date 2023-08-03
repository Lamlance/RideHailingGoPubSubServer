package ws

import (
	"goserver/libs"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func DriverHandlerMiddleware(c *fiber.Ctx) error {
	trip_id := c.Params("trip_id")
	if trip_id == "" {
		return c.SendStatus(400)
	}

	GlobalRoomMap.Lock.Lock()
	room, ok := GlobalRoomMap.Data[trip_id]
	if !ok {
		log.Println("Driver can't find trip Id in com channels")
		GlobalRoomMap.Lock.Unlock()
		return c.SendStatus(404)
	}
	log.Printf("Driver is stopping ride req loop")
	room.Break_ride_request <- true
	log.Printf("Driver has stopping ride req loop")
	c.Locals("room", room)

	GlobalRoomMap.Lock.Unlock()

	return c.Next()
}

func DriverListenThread(c *websocket.Conn) {
	defer c.Close()

	room, ok_room := c.Locals("room").(*CommunicationRoom)

	if !ok_room {
		log.Println("Driver cant find com channel")
		return
	}

	client_msg := room.client_msg
	driver_msg := room.driver_msg

	go DriverHandleClientMsgThread(c, client_msg)

	for {
		_, msg, err := c.ReadMessage()

		if err != nil {
			log.Println("Driver read error", err.Error())
			break
		} else {
			log.Println("Get driver msg: " + string(msg))
		}

		driver_msg.lock.Lock()
		driver_msg.data = libs.Enque(driver_msg.data, string(msg))
		driver_msg.lock.Unlock()
	}

}

func DriverHandleClientMsgThread(c *websocket.Conn, client_msg *CommunicationMsg) {
	for {
		client_msg.lock.Lock()

		if len(client_msg.data) != 0 {
			msg := ""
			msg, client_msg.data = libs.Dequeue(client_msg.data)
			log.Println("Driver get client msg: ", msg)

			err := c.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				client_msg.lock.Unlock()
				log.Println("Driver write error: " + err.Error())
				break
			}
		}
		client_msg.lock.Unlock()
		time.Sleep(2 * time.Second)
	}
}
