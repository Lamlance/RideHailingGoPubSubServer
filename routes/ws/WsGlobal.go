package ws

import "sync"

type CommunicationMsg struct {
	client_msg []string
	driver_msg []string

	lock *sync.Mutex
}

type GlobalCommunicationMsg struct {
	Data map[string]*CommunicationMsg
	Lock *sync.Mutex
}

var GlobalMsg = GlobalCommunicationMsg{
	Lock: new(sync.Mutex),
	Data: make(map[string]*CommunicationMsg),
}

func MakeEmptyCommunicationMsg() *CommunicationMsg {
	comMsg := CommunicationMsg{
		driver_msg: make([]string, 0),
		client_msg: make([]string, 0),
		lock:       new(sync.Mutex),
	}
	return &comMsg
}