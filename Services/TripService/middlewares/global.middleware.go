package middlewares

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

const (
	DriverFound      string = "DRF߷"
	NoDriver         string = "NDR߷"
	DriverCancel     string = "DCX߷"
	DriverArrivePick string = "DAP߷"
	DriverStratTrip  string = "DST߷"
	DriverArriveDrop string = "DAD߷"
	ClientCancel     string = "CCX߷"
	TripId           string = "TID߷"
	Message          string = "MSG߷"
)

var GlobalRoomMap = GlobalCommunicationMsg{
	Lock: new(sync.Mutex),
	Data: make(map[string]*CommunicationRoom),
}

func MakeEmptyCommunicationRoom() *CommunicationRoom {
	comMsg := CommunicationRoom{
		driver_msg: &CommunicationMsg{
			data: make([]string, 0),
			lock: new(sync.Mutex),
		},
		client_msg: &CommunicationMsg{
			data: make([]string, 0),
			lock: new(sync.Mutex),
		},
		Ride_requst_channel: make(chan int, 0),
		lock:                new(sync.Mutex),
	}
	return &comMsg
}

func GetSocketMsgHandler(msg string, c *websocket.Conn) error {
	var err error = nil
	switch msg[0:5] {
	case NoDriver:
		err = c.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(3001, "No driver found"))
	case ClientCancel:
		err = c.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(3001, "Client has canceled trip"))
	case DriverCancel:
		err = c.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(3002, "Driver has canceled trip"))
	case DriverArriveDrop:
		err = c.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(3003, "Trip has finished"))
	}
	return err
}

func RecevideSocketMsgHandler(msg string, c *websocket.Conn) error {
	var err error = nil
	switch msg[0:5] {
	case NoDriver:
		err = c.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(3001, NoDriver+"No driver found"))
	case ClientCancel:
		err = c.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(3001, ClientCancel+"Client has canceled trip"))
	case DriverCancel:
		err = c.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(3002, DriverCancel+"Driver has canceled trip"))
	case DriverArriveDrop:
		err = c.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(3003, DriverArriveDrop+"Trip has finished"))
	case Message, DriverFound, DriverArrivePick, DriverStratTrip:
		err = c.WriteMessage(websocket.TextMessage, []byte(msg))
	default:
		err = c.WriteMessage(websocket.TextMessage, []byte(Message+msg))
	}

	return err
}

func Enque(queue []string, element string) []string {
	queue = append(queue, element) // Simply append to enqueue.
	return queue
}

func Dequeue(queue []string) (string, []string) {
	first := queue[0]
	if len(queue) == 1 {
		var tmp = []string{}
		return first, tmp
	}

	return first, queue[1:]
}
