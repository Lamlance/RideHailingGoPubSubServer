package ws

import (
	"encoding/json"
	"goserver/libs"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func DriverHandlerMiddleware(c *fiber.Ctx) error {
	trip_id := c.Params("trip_id")
	driver_id := c.Query("driver_id")
	if trip_id == "" || driver_id == "" {
		return c.SendStatus(400)
	}

	GlobalRoomMap.Lock.Lock()
	room, ok := GlobalRoomMap.Data[trip_id]
	if !ok {
		log.Println("Driver can't find trip Id in com channels")
		GlobalRoomMap.Lock.Unlock()
		return c.SendStatus(404)
	}
	c.Locals("room", room)

	GlobalRoomMap.Lock.Unlock()

	return c.Next()
}

func DriverListenThread(c *websocket.Conn) {
	room, ok_room := c.Locals("room").(*CommunicationRoom)

	if !ok_room {
		log.Println("Driver cant find com channel")
		c.Close()
		return
	}

	log.Printf("Driver is stopping ride req loop")
	room.Ride_requst_channel <- 0
	log.Printf("Driver has stopping ride req loop")

	client_msg := room.client_msg
	driver_msg := room.driver_msg

	go DriverHandleClientMsgThread(c, client_msg)

	driver_info := struct {
		Driver_id string `json:"driver_id"`
	}{
		Driver_id: c.Query("driver_id"),
	}
	data, _ := json.Marshal(driver_info)
	driver_msg.lock.Lock()
	driver_msg.data = libs.Enque(driver_msg.data, DriverFound+string(data))
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

		if msg[0:5] == NoDriver || msg[0:5] == ClientCancel || msg[0:5] == DriverCancel {
			running = false
		}

		driver_msg.lock.Lock()
		driver_msg.data = libs.Enque(driver_msg.data, msg)
		driver_msg.lock.Unlock()
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
		msg, client_msg.data = libs.Dequeue(client_msg.data)
		log.Println("Driver get client msg: ", msg)

		var err error

		switch msg[0:5] {
		case NoDriver:
			running = false
		case ClientCancel:
			err = c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(3001, "Client has canceled trip"))
		case DriverCancel:
			err = c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(3002, "Driver has canceled trip"))
		case Message:
			err = c.WriteMessage(websocket.TextMessage, []byte(msg))
		case DriverFound:
			err = c.WriteMessage(websocket.TextMessage, []byte(msg))
		default:
			err = c.WriteMessage(websocket.TextMessage, []byte(Message+msg))
		}

		if err != nil {
			client_msg.lock.Unlock()
			log.Println("Driver write error: " + err.Error())
			break
		}

		client_msg.lock.Unlock()
	}
	c.Close()
}
