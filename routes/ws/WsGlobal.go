package ws

import (
	"sync"
)
type RideReqInfo struct {
	SLon float64 `json:"slon"`
	SLat float64 `json:"slat"`
	SAdr string  `json:"sadr"`

	ELon float64 `json:"elon"`
	ELat float64 `json:"elat"`
	EAdr string `json:"eadr"`

	User_id string `json:"user_id"`
	Driver_id string `json:"driver_id"`
	Trip_id string `json:"trip_id"`
}

type CommunicationMsg struct {
	data []string
	lock *sync.Mutex
}

type CommunicationRoom struct {
	client_msg *CommunicationMsg
	driver_msg *CommunicationMsg

	RideInfo *RideReqInfo

	lock               *sync.Mutex 
	Ride_requst_channel chan int
}

type GlobalCommunicationMsg struct {
	Data map[string]*CommunicationRoom
	Lock *sync.Mutex
}



const (
	DriverFound string = "⚼DRF"
	NoDriver string = "⚼NDR"
	DriverCancel string = "⚼DCX"
	ClientCancel string = "⚼CCX"
	TripId string = "⚼TID"
	Message string = "⚼MSG"
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
		Ride_requst_channel: make(chan int,0),
		lock: new(sync.Mutex),
	}
	return &comMsg
}
