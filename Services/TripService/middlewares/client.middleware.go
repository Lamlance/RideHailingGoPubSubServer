package middlewares

import (
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

func ClientListenThread(c *websocket.Conn) {
	trip_id, ok_parseId := c.Locals("trip_id").(string)
	room, ok_room := c.Locals("room").(*CommunicationRoom)

	if !ok_parseId || !ok_room {
		log.Println("Locals error", ok_parseId, ok_room)
		c.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(3000, "Server error"))
		return
	}

	client_msg := room.client_msg
	driver_msg := room.driver_msg

	c.SetCloseHandler(func(code int, text string) error {
		log.Println("Finishing trip: ", trip_id)
		GlobalRoomMap.Lock.Lock()
		_, ok := GlobalRoomMap.Data[trip_id]
		if ok {
			delete(GlobalRoomMap.Data, trip_id)
		}
		GlobalRoomMap.Lock.Unlock()
		return nil
	})

	go ClientHandleDriverMsgThread(c,driver_msg)

	running := true
	for running {
		_, data, err := c.ReadMessage()
		if err != nil {
			log.Println("Client read error: " + err.Error())
			break
		}
		msg := string(data)
		log.Println("Get client msg: " + msg)

		client_msg.lock.Lock()
		client_msg.data = Enque(client_msg.data, msg)
		client_msg.lock.Unlock()

		err = GetSocketMsgHandler(msg, c)
		if err != nil {
			log.Println("Client handle msg error: ", err)
			break
		}

	}
}

func ClientHandleDriverMsgThread(c *websocket.Conn, driver_msg *CommunicationMsg) {
	running := true
	for ; running; time.Sleep(2 * time.Second) {
		driver_msg.lock.Lock()
		if len(driver_msg.data) <= 0 {
			driver_msg.lock.Unlock()
			continue
		}

		msg := ""
		msg, driver_msg.data = Dequeue(driver_msg.data)
		log.Println("Client get msg: ", msg)
		var err error

		err = RecevideSocketMsgHandler(msg, c)

		if err != nil {
			driver_msg.lock.Unlock()
			log.Println("Client write error: " + err.Error())
			break
		}

		driver_msg.lock.Unlock()

	}
	c.Close()
}
