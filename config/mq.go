package config

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
) 

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func SetupMQ(queueName string) *RabbitMQ {

	// Load url from .env
	url := os.Getenv("AMQP_URL")

	// Connect to MQ
	log.Println("Connecting to MQ")
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Println("Connected to MQ!")

	// Open Channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	// Create Queue
	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")
	log.Printf("Queue %s declared!\n", queueName)

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
		Queue:   q,
	}
}