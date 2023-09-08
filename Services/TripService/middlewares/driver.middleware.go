package middlewares

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

func DriverListenThread(c *websocket.Conn) {
	room, ok_room := c.Locals("room").(*CommunicationRoom)
	trip_id := c.Params("trip_id")

	if !ok_room {
		log.Println("Driver cant find com channel")
		c.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(3000, "Server error"))
		return
	}

	c.SetCloseHandler(func(code int, text string) error {
		log.Println("Finishing trip: ", trip_id)
		GlobalRoomMap.Lock.Lock()
		_, ok := GlobalRoomMap.Data[trip_id]
		if ok && code >= 3000 {
			delete(GlobalRoomMap.Data, trip_id)
		}
		GlobalRoomMap.Lock.Unlock()
		return nil
	})

	client_msg := room.client_msg
	driver_msg := room.driver_msg

	go DriverHandleClientMsgThread(c, client_msg)

	room.RideInfo.Driver_id = c.Query("driver_id")

	data, _ := json.Marshal(room.RideInfo)
	driver_msg.lock.Lock()
	driver_msg.data = Enque(driver_msg.data, DriverFound+string(data))
	driver_msg.lock.Unlock()

	running := true
	for running {
		_, data, err := c.ReadMessage()
		if err != nil {
			log.Println("Driver read error", err.Error())
			break
		}

		msg := string(data)
		log.Println("Get driver msg: " + msg)

		driver_msg.lock.Lock()
		driver_msg.data = Enque(driver_msg.data, msg)
		driver_msg.lock.Unlock()

		err = GetSocketMsgHandler(msg, c)
		if err != nil {
			log.Println("Client handle msg error: ", err)
			break
		}
	}
}


func DriverHandleClientMsgThread(c *websocket.Conn, client_msg *CommunicationMsg) {
	running := true
	for ; running; time.Sleep(2 * time.Second) {
		client_msg.lock.Lock()
		if len(client_msg.data) <= 0 {
			client_msg.lock.Unlock()
			continue
		}

		msg := ""
		msg, client_msg.data = Dequeue(client_msg.data)
		log.Println("Driver get client msg: ", msg)

		var err error

		err = RecevideSocketMsgHandler(msg, c)

		if err != nil {
			client_msg.lock.Unlock()
			log.Println("Driver write error: " + err.Error())
			break
		}

		client_msg.lock.Unlock()
	}
	c.Close()
}
