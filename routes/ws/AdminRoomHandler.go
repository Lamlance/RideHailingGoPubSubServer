package ws

import (
	"encoding/json"
	"goserver/libs"
	"goserver/routes"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

func AdminRoomHandler(c *websocket.Conn) {
	_, ok_parseId := c.Locals("trip_id").(string)
	room, ok_room := c.Locals("room").(*CommunicationRoom)

	if !ok_parseId || !ok_room {
		log.Println("Locals error", ok_parseId, ok_room)
		c.Close()
		return
	}

	driver_msg := room.driver_msg
	running := true
	for ; running; time.Sleep(2 * time.Second) {
		driver_msg.lock.Lock()
		if len(driver_msg.data) <= 0 {
			driver_msg.lock.Unlock()
			continue
		}

		msg := ""
		msg, driver_msg.data = libs.Dequeue(driver_msg.data)
		log.Println("Admin get msg: ", msg)
		var err error

		switch msg[0:5] {
		case NoDriver:
			running = false
			err = c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(3000, NoDriver+"No driver found"))
		case ClientCancel:
			running = false
			err = c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(3001, ClientCancel+"Client has canceled trip"))
		case DriverCancel:
			running = false
			err = c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(3002, DriverCancel+"Driver has canceled trip"))
		case DriverFound:
			info := &routes.RideReqInfo{}
			json.Unmarshal([]byte(msg[5:]), info)

			new_info := &struct {
				Driver_id string `json:"driver_id"`
			}{
				Driver_id: info.Driver_id,
			}
			new_msg, _ := json.Marshal(new_info)

			err = c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(3999, DriverFound+string(new_msg)))
		}

		if err != nil {
			break
		}
	}

	c.Close()
}
