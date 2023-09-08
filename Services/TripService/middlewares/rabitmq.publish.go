package middlewares

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RideReqToPub = make(chan *RideReqInfo, 50)

func PublishRideRequest() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panic("Listen Ride Req error:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panic("Create channel error:", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"create_ride_req_topic", // name
		"topic",                 // type
		true,                    // durable
		false,                   // auto-deleted
		false,                   // internal
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		log.Panic("Exchange error:", err)
	}
	ctx := context.Background()

	for {
		data, ok := <-RideReqToPub
		if !ok {
			log.Println("Rabbit mq cant get ride req to pub")
			continue
		}
		b, err := json.Marshal(data)
		if err != nil {
			log.Println("Ride req to pub json error: ", err)
			continue
		}

		err = ch.PublishWithContext(ctx,
			"create_ride_req_topic", // exchange
			"anonymous.info",        // routing key
			false,                   // mandatory
			false,                    // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        b,
			},
		)
		if err != nil {
			log.Println("Rabbit mq publish error: ", err)
		} else {
			log.Println("Rabbitmq published trip req")
		}
	}

}
