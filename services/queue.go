package services

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/config"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/util"
)

type Callback func(ch *amqp.Channel, ctx context.Context, msg []byte, routingKey string)

func Consume(rmq *config.RabbitMQ, queueName string, callback Callback) {
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
			ctx := context.Background()

			// Process the message
			log.Printf("Received a message from queue %s: %s", queueName, d.Body)

			// Get Routing Key
			routingKey := util.DetermineNewRoutingKey(d.RoutingKey)

			// Process message
			callback(ch, ctx, d.Body, routingKey)

			// Acknowledge the message
			err = d.Ack(false)
			util.FailOnError(err, "Failed to ack")
		}
	}()

	<-forever
}

func Publish(ch *amqp.Channel, ctx context.Context, msg []byte, routingKey string) {
	err := ch.PublishWithContext(
		ctx,
		"topic.router", // exchange
		routingKey,     // routing key
		true,           // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	util.FailOnError(err, "Failed to publish a message")
	log.Printf("Published a message with routing key %s", routingKey)

}
