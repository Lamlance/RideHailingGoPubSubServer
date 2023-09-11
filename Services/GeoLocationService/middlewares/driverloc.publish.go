package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

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
		log.Panicf(msg+" ", err)
	}
}
var RabbitMQCon *amqp.Connection = nil
func ConnectToRabbitMQ() {
	port := os.Getenv("RABBITMQ_PORT")
	host := os.Getenv("RABBITMQ_HOST")

	if port == "" {
		port = "5672"
	}
	if host == "" {
		host = "localhost"
	}

	url := fmt.Sprintf("amqp://guest:guest@%s:%s", host, port)
	log.Println("Read rabbit mq link: ", url)

	conn, err := amqp.Dial(url)
	if err != nil {
		failOnError(err, "Connection error")
	}
	RabbitMQCon = conn
}

func PublishDriverLocation() {
	if RabbitMQCon == nil {
		log.Panicln("Cant connect to rabbit mq")
	}
	conn := RabbitMQCon
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
