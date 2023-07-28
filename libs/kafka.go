package libs

import (
	"log"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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

	return &ps
}

// var GlobalKafkaConsumer *kafka.Consumer
var GlobalPubSubDict = CreatedPubSubDict{
	pubsubs:map[string]*PubSub{},
	lock: new(sync.Mutex),
}

func Subscribe(p *PubSub, topic string) (<-chan string, func()) {
	p.lock.Lock()
	defer p.lock.Unlock()

	c := make(chan string, 1)
	p.channels = append(p.channels, c)

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
	}
}

func Publish(topic string, message string) {
	ps, ok := GlobalPubSubDict.pubsubs[topic]
	if !ok {
		return
	}

	for _, c := range ps.channels {
		c <- message
	}
}

func KafkaConsumer() {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "DriverConsumer",
	})

	if err != nil {
		log.Println("Kafka connect err: ",err)
		return
	}

	log.Println("Kafka consumer: ", c)

	err = c.Subscribe("Driver", nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

	run := true
	for run {
		msg, err := c.ReadMessage(10 * time.Minute)

		if err != nil {
			run = err.(kafka.Error).Code() == kafka.ErrTimedOut
			if run {
				log.Println("Kafka read message timeout")
			} else {
				log.Println("Kafka read error: ", err)
			}
		}

		log.Println("Consumer message: ", string(msg.Value))
		Publish("Driver", string(msg.Value))
	}

	log.Println("Closing Consumer")
}
