package middlewares

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type DriverPubLoc struct {
	Lon       float64 `json:"lon"`
	Lat       float64 `json:"lat"`
	Driver_id string  `json:"driver_id"`
}

var DriverLocationToPub = make(chan *DriverPubLoc, 50)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func PublishDriverLocation() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Listen Ride Req error")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Create channel error")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"update_driver_loc", // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Exchange declare error")

	ctx := context.Background()
	for {
		data, ok := <-DriverLocationToPub
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
			"update_driver_loc", // exchange
			"anonymous.info",    // routing key
			false,               // mandatory
			false,               // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        b,
			},
		)
		if err != nil {
			log.Println("Rabbit mq publish error: ", err)
		} else {
			log.Println("Rabbitmq published driver coord: ", string(b))
		}
	}
}
