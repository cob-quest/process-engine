package config

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/util"
)

type RabbitMQ struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func SetupMQ() *RabbitMQ {

	// amqp
	AMQP_PORT := os.Getenv("AMQP_PORT")
	AMQP_HOSTNAME := os.Getenv("AMQP_HOSTNAME")
	AMQP_USERNAME := os.Getenv("AMQP_USERNAME")
	AMQP_PASSWORD := os.Getenv("AMQP_PASSWORD")

	// Load url from .env
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		AMQP_USERNAME,
		AMQP_PASSWORD,
		AMQP_HOSTNAME,
		AMQP_PORT,
	)

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
		"queue.trigger.fromService", // name
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")
	log.Printf("%s created!\n", q.Name)

	// Notification From Queue
	q2, err := ch.QueueDeclare(
		"queue.notification.fromService", // name
		true,                             // durable
		false,                            // delete when unused
		false,                            // exclusive
		false,                            // no-wait
		nil,                              // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")
	log.Printf("%s created!\n", q2.Name)

	// Notification To Queue
	q3, err := ch.QueueDeclare(
		"queue.notification.toService", // name
		true,                           // durable
		false,                          // delete when unused
		false,                          // exclusive
		false,                          // no-wait
		nil,                            // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")
	log.Printf("%s created!\n", q3.Name)

	// Image Builder From Queue
	q4, err := ch.QueueDeclare(
		"queue.imageBuilder.fromService", // name
		true,                             // durable
		false,                            // delete when unused
		false,                            // exclusive
		false,                            // no-wait
		nil,                              // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")
	log.Printf("%s created!\n", q4.Name)

	// Image Builder To Queue
	q5, err := ch.QueueDeclare(
		"queue.imageBuilder.toService", // name
		true,                           // durable
		false,                          // delete when unused
		false,                          // exclusive
		false,                          // no-wait
		nil,                            // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")
	log.Printf("%s created!\n", q5.Name)

	// challenge From Queue
	q6, err := ch.QueueDeclare(
		"queue.challenge.fromService", // name
		true,                           // durable
		false,                          // delete when unused
		false,                          // exclusive
		false,                          // no-wait
		nil,                            // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")
	log.Printf("%s created!\n", q6.Name)

	// challenge To Queue
	q7, err := ch.QueueDeclare(
		"queue.challenge.toService", // name
		true,                         // durable
		false,                        // delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil,                          // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")
	log.Printf("%s created!\n", q7.Name)

	// ------------------
	// EXCHANGE CREATION
	// ------------------

	// Router Topic
	err = ch.ExchangeDeclare(
		"topic.router", // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")
	log.Println("topic.router created!")

	// Trigger Topic
	err = ch.ExchangeDeclare(
		"topic.trigger", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")
	log.Println("topic.trigger created!")

	// Notification Topic
	err = ch.ExchangeDeclare(
		"topic.notification", // name
		"topic",              // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")
	log.Println("topic.notification created!")

	// challenge Topic
	err = ch.ExchangeDeclare(
		"topic.challenge", // name
		"topic",            // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")
	log.Println("topic.challenge created!")

	// Image Builder Topic
	err = ch.ExchangeDeclare(
		"topic.imageBuilder", // name
		"topic",              // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
	util.FailOnError(err, "Failed to declare an exchange")
	log.Println("topic.imageBuilder created!")

	// -----------------
	// BINDING CREATION
	// -----------------

	// Router Topic - Image Builder Topic
	err = ch.ExchangeBind(
		"topic.imageBuilder", // exchange destination name
		"imageBuilder.#",     // routing key
		"topic.router",  // exchange source name
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")
	log.Println("Router Topic - Image Builder Topic")

	// Router Topic - Notification Topic
	err = ch.ExchangeBind(
		"topic.notification", // exchange destination name
		"notification.#",     // routing key
		"topic.router",       // exchange source name
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")
	log.Println("Router Topic - Notification Topic")

	// Router Topic - challenge Topic
	err = ch.ExchangeBind(
		"topic.challenge", // exchange destination name
		"challenge.#",     // routing key
		"topic.router",     // exchange source name
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")
	log.Println("Router Topic - challenge Topic")

	// Trigger Topic - Trigger FromQueue
	err = ch.QueueBind(
		"queue.trigger.fromService", // queue name
		"trigger.fromService.*",     // routing key
		"topic.trigger",             // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")
	log.Println("Trigger Topic - Trigger FromQueue")

	// Notification Topic - Notification FromQueue
	err = ch.QueueBind(
		"queue.notification.fromService", // queue name
		"notification.fromService.*",     // routing key
		"topic.notification",             // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")
	log.Println("Notification Topic - Notification FromQueue")

	// Notification Topic - Notification ToQueue
	err = ch.QueueBind(
		"queue.notification.toService", // queue name
		"notification.toService.*",     // routing key
		"topic.notification",           // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")
	log.Println("Notification Topic - Notification ToQueue")

	// Image Builder Topic - Image Builder FromQueue
	err = ch.QueueBind(
		"queue.imageBuilder.fromService", // queue name
		"imageBuilder.fromService.*",     // routing key
		"topic.imageBuilder",             // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")
	log.Println("Image Builder Topic - Image Builder FromQueue")

	// Image Builder Topic - Image Builder ToQueue
	err = ch.QueueBind(
		"queue.imageBuilder.toService", // queue name
		"imageBuilder.toService.*",     // routing key
		"topic.imageBuilder",           // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")
	log.Println("Image Builder Topic - Image Builder ToQueue")

	// challenge Topic - challenge FromQueue
	err = ch.QueueBind(
		"queue.challenge.fromService", // queue name
		"challenge.fromService.*",     // routing key
		"topic.challenge",             // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")
	log.Println("challenge Topic - challenge FromQueue")

	// challenge Topic - challenge ToQueue
	err = ch.QueueBind(
		"queue.challenge.toService", // queue name
		"challenge.toService.*",     // routing key
		"topic.challenge",           // exchange
		false,
		nil,
	)
	util.FailOnError(err, "Failed to bind a queue")
	log.Println("challenge Topic - challenge ToQueue")

	return &RabbitMQ{
		Conn: conn,
		Ch:   ch,
	}
}
