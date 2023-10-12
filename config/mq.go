package config

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/util"
) 

type RabbitMQ struct {
	Conn    *amqp.Connection
	Ch *amqp.Channel
}

func SetupMQ() *RabbitMQ {

	// Load url from .env
	url := os.Getenv("AMQP_URL")

	// Connect to MQ
	log.Println("Connecting to MQ")
	conn, err := amqp.Dial(url)
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	log.Println("Connected to MQ!")

	// Open Channel
	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")

	// ---------------
	// QUEUE CREATION
	// ---------------

	// Trigger From Queue
	q, err := ch.QueueDeclare(
		"queue.trigger.fromService",    // name
		true, // durable
		false, // delete when unused
		false,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")
	log.Printf("%s created!\n", q.Name)

	// Notification From Queue
	q2, err := ch.QueueDeclare(
		"queue.notification.fromService",    // name
		true, // durable
		false, // delete when unused
		false,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")
	log.Printf("%s created!\n", q2.Name)

	// Notification To Queue
	q3, err := ch.QueueDeclare(
		"queue.notification.toService",    // name
		true, // durable
		false, // delete when unused
		false,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")
	log.Printf("%s created!\n", q3.Name)

	// ------------------
	// EXCHANGE CREATION
	// ------------------

	// Trigger Topic
	err = ch.ExchangeDeclare(
		"topic.trigger", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")
	log.Println("topic.trigger created!")

	// Notification Topic
	err = ch.ExchangeDeclare(
		"topic.notification", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")
	log.Println("topic.notification created!")

	// -----------------
	// BINDING CREATION
	// -----------------

	// Trigger Topic - Trigger FromQueue
	err = ch.QueueBind(
		"queue.trigger.fromService",       // queue name
		"trigger.fromService.*",           // routing key
		"topic.trigger",       // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")
	log.Println("topic.trigger bound to queue.trigger.fromService with routing key: trigger.*")

	// Notification Topic - Notification ToQueue
	err = ch.QueueBind(
		"queue.notification.toService",       // queue name
		"notification.toService.*",           // routing key
		"topic.notification",       // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")
	log.Println("Notification Topic - Notification ToQueue")

	// Notification Topic - Notification FromQueue
	err = ch.QueueBind(
		"queue.notification.fromService",       // queue name
		"notification.fromService.*",           // routing key
		"topic.notification",       // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")
	log.Println("Notification Topic - Notification FromQueue")

	return &RabbitMQ{
		Conn:    conn,
		Ch: ch,
	}
}