package ws

import "sync"

type CommunicationMsg struct {
	data []string
	lock *sync.Mutex
}

type CommunicationRoom struct {
	client_msg *CommunicationMsg
	driver_msg *CommunicationMsg

	lock               *sync.Mutex 
	Break_ride_request chan bool
}

type GlobalCommunicationMsg struct {
	Data map[string]*CommunicationRoom
	Lock *sync.Mutex
}

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
		Break_ride_request: make(chan bool,0),
		lock: new(sync.Mutex),
	}
	return &comMsg
}
