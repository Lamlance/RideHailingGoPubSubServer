package ws

import (
	"goserver/libs"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func ClientCheckMiddleware(c *fiber.Ctx) error {
	random_string := "Secrete_trip_id" //uuid.New().String()

	GlobalMsg.Lock.Lock()
	GlobalMsg.Data[random_string] = MakeEmptyCommunicationMsg()
	GlobalMsg.Lock.Unlock()

	c.Locals("trip_id", random_string)
	return c.Next()
}

func ClientListenThread(c *websocket.Conn) {
	trip_id, ok_parseId := c.Locals("trip_id").(string)

	if !ok_parseId {
		log.Println("Can't parse trip id to string")
		c.Close()
		return
	}

	log.Println("Client trip id: ", trip_id)
	if trip_id == "" {
		c.Close()
		return
	}

	GlobalMsg.Lock.Lock()
	comMsg, ok := GlobalMsg.Data[trip_id]
	GlobalMsg.Lock.Unlock()

	if !ok {
		c.Close()
		log.Println("Client can't find trip Id: " + trip_id)
		return
	}

	go ClientHandleDriverMsgThread(c, comMsg)

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Client read error: " + err.Error())
			break
		} else {
			log.Println("Get client msg: " + string(msg))
		}

		comMsg.lock.Lock()
		comMsg.client_msg = libs.Enque(comMsg.client_msg, string(msg))
		comMsg.lock.Unlock()
	}

	c.Close()
}

func ClientHandleDriverMsgThread(c *websocket.Conn, comMsg *CommunicationMsg) {
	for {
		comMsg.lock.Lock()

		if len(comMsg.driver_msg) != 0 {
			msg := ""
			msg, comMsg.driver_msg = libs.Dequeue(comMsg.driver_msg)
			log.Println("Client get driver msg: ", msg)

			err := c.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				comMsg.lock.Unlock()
				log.Println("Client write error: " + err.Error())
				break
			}
		}
		comMsg.lock.Unlock()
		time.Sleep(2 * time.Second)

	}
}