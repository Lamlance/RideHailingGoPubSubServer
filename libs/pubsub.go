package libs

import (
	"log"
	"strconv"
	"sync"
)

type KafkaMsg struct {
	msg     string
	created int64
}

type PubSub struct {
	channels []chan string
	lock     *sync.RWMutex
}

type CreatedPubSubDict struct {
	pubsubs map[string]*PubSub
	lock    *sync.Mutex
}

func NewPubSub(topic string) *PubSub {
	ps := PubSub{
		channels: make([]chan string, 0),
		lock:     new(sync.RWMutex),
	}

	GlobalPubSubDict.lock.Lock()
	GlobalPubSubDict.pubsubs[topic] = &ps
	GlobalPubSubDict.lock.Unlock()

	log.Println("Create topic: ",topic)

	return &ps
}

// var GlobalKafkaConsumer *kafka.Consumer
var GlobalPubSubDict = CreatedPubSubDict{
	pubsubs:map[string]*PubSub{},
	lock: new(sync.Mutex),
}

func Subscribe(topic string) (<-chan string, func(),bool) {
	GlobalPubSubDict.lock.Lock()
	defer GlobalPubSubDict.lock.Unlock()
	p,ok := GlobalPubSubDict.pubsubs[topic]
	if !ok{
		log.Println("Cant find topic ",topic, " in pubsub dict")
		return nil,nil, false
	}
	
	p.lock.Lock()
	defer p.lock.Unlock()

	c := make(chan string, 10)
	p.channels = append(p.channels, c)
	log.Println("Channel length: ",topic,len(p.channels))

	return c, func() {
		p.lock.Lock()
		defer p.lock.Unlock()
		for i, channel := range p.channels {
			if channel == c {
				p.channels = append(p.channels[:i], p.channels[i+1:]...)
				close(c)
				return
			}
		}
	},true
}

func Publish(topic string, message string) {
	GlobalPubSubDict.lock.Lock()
	defer GlobalPubSubDict.lock.Unlock()

	ps, ok := GlobalPubSubDict.pubsubs[topic]
	if !ok {
		return
	}

	for i, c := range ps.channels {
		log.Println("Pub topic: ",topic,strconv.Itoa(i))
		c <- message
	}
}
