package ws

import (
	"goserver/routes"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type CommunicationMsg struct {
	data []string
	lock *sync.Mutex
}

type CommunicationRoom struct {
	client_msg *CommunicationMsg
	driver_msg *CommunicationMsg

	RideInfo *routes.RideReqInfo

	lock                *sync.Mutex
	Ride_requst_channel chan int
}

type GlobalCommunicationMsg struct {
	Data map[string]*CommunicationRoom
	Lock *sync.Mutex
}

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


func RecevideSocketMsgHandler(msg string, c *websocket.Conn) error {
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
	case Message, DriverFound, DriverArriveDrop, DriverArrivePick,DriverStratTrip:
		err = c.WriteMessage(websocket.TextMessage, []byte(msg))
	default:
		err = c.WriteMessage(websocket.TextMessage, []byte(Message+msg))
	}

	return err
}
