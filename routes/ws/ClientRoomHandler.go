package ws

import (
	"goserver/libs"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

func ClientListenThread(c *websocket.Conn) {
	defer c.Close()

	_, ok_parseId := c.Locals("trip_id").(string)
	room, ok_room := c.Locals("room").(*CommunicationRoom)

	if !ok_parseId || !ok_room{
		log.Println("Locals error",ok_parseId,ok_room)
		c.Close()
		return
	}

	client_msg := room.client_msg
	driver_msg := room.driver_msg

	go ClientHandleDriverMsgThread(c, driver_msg)

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Client read error: " + err.Error())
			break
		} else {
			log.Println("Get client msg: " + string(msg))
		}

		client_msg.lock.Lock()
		client_msg.data = libs.Enque(client_msg.data, string(msg))
		client_msg.lock.Unlock()
	}

}

func ClientHandleDriverMsgThread(c *websocket.Conn, driver_msg *CommunicationMsg) {
	for {
		driver_msg.lock.Lock()

		if len(driver_msg.data) != 0 {
			msg := ""
			msg, driver_msg.data = libs.Dequeue(driver_msg.data)
			log.Println("Client get driver msg: ", msg)

			err := c.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				driver_msg.lock.Unlock()
				log.Println("Client write error: " + err.Error())
				break
			}
		}
		driver_msg.lock.Unlock()
		time.Sleep(2 * time.Second)

	}
}
