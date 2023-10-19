package main

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/config"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/util"
)

func main() {
	
	// // Load env file
	// log.Println("Loading .env file")
	// config.InitEnvironment()
	// log.Println(".env loaded!")

	rmq := config.SetupMQ()
	defer rmq.Conn.Close()
	defer rmq.Ch.Close()

	go consume(rmq, "queue.notification.fromService")
	go consume(rmq, "queue.trigger.fromService")
	go consume(rmq, "queue.assignment.fromService")
	go consume(rmq, "queue.imageBuilder.fromService")

	select {}
	// msgs, err := rmq.Ch.Consume(
	// 	"queue.trigger.fromService", // queue
	// 	"processEngine",             // consumer
	// 	false,                       // auto-ack
	// 	false,                       // exclusive
	// 	false,                       // no-local
	// 	false,                       // no-wait
	// 	nil,                         // args
	// )
	// util.FailOnError(err, "Failed to consume messages from queue.trigger.fromService")

	// var forever chan struct{}

	// go func() {
	// 	log.Println("Consuming")
	// 	for msg := range msgs {
	// 		log.Printf("Received msg: %s\n", msg.Body)

	// 		err = msg.Ack(false)
	// 		if err != nil {
	// 			log.Printf("Failed to ack message: %s", err)
	// 		}
	// 	}
	// }()

	// log.Printf(" [*] Waiting for messages")
	// <-forever
}

func consume(rmq *config.RabbitMQ, queueName string) {
	// Create a new channel for this queue
	ch, err := rmq.Conn.Channel()
	util.FailOnError(err, "Failed to open a channel")

	defer ch.Close()

	msgs, err := ch.Consume(
		queueName,
		"processEngine", // consumer
		false,           // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	util.FailOnError(err, fmt.Sprintf("Failed to register a consumer for queue %s: %s", queueName, err))

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Process the message
			log.Printf("Received a message from queue %s: %s", queueName, d.Body)

			// After processing, publish to another queue
			ctx := context.Background()

			body := "Processed message from " + queueName + ": " + string(d.Body)
			err := ch.PublishWithContext(
				ctx,
				"topic.notification",               // exchange
				"notification.toService.testEvent", // routing key
				true,                               // mandatory
				false,                              // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			util.FailOnError(err, "Failed to publish a message")

			// Acknowledge the message
			d.Ack(false)
		}
	}()

	<-forever
}
