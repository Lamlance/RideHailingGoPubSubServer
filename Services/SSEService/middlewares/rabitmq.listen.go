package middlewares

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RideReqInfo struct {
	SLon float64 `json:"slon"`
	SLat float64 `json:"slat"`
	SAdr string  `json:"sadr"`

	ELon float64 `json:"elon"`
	ELat float64 `json:"elat"`
	EAdr string  `json:"eadr"`

	User_id   string `json:"user_id"`
	Driver_id string `json:"driver_id"`
	Trip_id   string `json:"trip_id"`

	Price float64 `json:"price"`

	User_Name  string `json:"user_name"`
	User_Phone string `json:"user_phone"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var RabbitMQCon1 *amqp.Connection = nil
var RabbitMQCon2 *amqp.Connection = nil

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

	conn1, err := amqp.Dial(url)
	conn2, err := amqp.Dial(url)

	if err != nil {
		failOnError(err, "Connection error")
	}
	RabbitMQCon1 = conn1
	RabbitMQCon2 = conn2
}

func ListenDriverLoc() {
	conn := RabbitMQCon1
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
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
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,              // queue name
		"anonymous.info",    // routing key
		"update_driver_loc", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	for {
		data, ok := <-msgs
		if !ok {
			log.Println("Consume message error")
			continue
		} else {
			log.Println("Get watch coord ")
		}
		Publish("DriverCoord", string(data.Body))
	}
}

func ListenRideRequest() {
	conn := RabbitMQCon2
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
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
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,                  // queue name
		"anonymous.info",        // routing key
		"create_ride_req_topic", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	for {
		data, ok := <-msgs
		if !ok {
			log.Println("Consume message error")
			continue
		} else {
			log.Println("Get ride req ")
		}
		Publish("w3gv", string(data.Body))
	}
}
