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
		c.SendStatus(400)
	}

	GlobalMsg.Lock.Lock()
	comMsg, ok := GlobalMsg.Data[trip_id]
	GlobalMsg.Lock.Unlock()

	if !ok {
		log.Println("Driver can't find trip Id in com channels")
		return c.SendStatus(404)
	}

	c.Locals("ComChannel",comMsg)

	return c.Next()
}

func DriverListenThread(c *websocket.Conn) {
	defer c.Close()

	comMsg,ok := c.Locals("ComChannel").(*CommunicationMsg)
	if !ok{
		log.Println("Driver cant find com channel")
		return
	}


	go DriverHandleClientMsgThread(c, comMsg)

	for {
		_, msg, err := c.ReadMessage()

		if err != nil {
			log.Println("Driver read error", err.Error())
			break
		} else {
			log.Println("Get driver msg: " + string(msg))
		}

		comMsg.lock.Lock()
		comMsg.driver_msg = libs.Enque(comMsg.driver_msg, string(msg))
		comMsg.lock.Unlock()
	}

	c.Close()
}

func DriverHandleClientMsgThread(c *websocket.Conn, comMsg *CommunicationMsg) {
	for {
		comMsg.lock.Lock()

		if len(comMsg.client_msg) != 0 {
			msg := ""
			msg, comMsg.client_msg = libs.Dequeue(comMsg.client_msg)
			log.Println("Driver get client msg: ", msg)

			err := c.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				comMsg.lock.Unlock()
				log.Println("Driver write error: " + err.Error())
				break
			}
		}
		comMsg.lock.Unlock()
		time.Sleep(2 * time.Second)

	}
}