package services

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/config"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/util"
)

type Callback func(ch *amqp.Channel, ctx context.Context, d amqp.Delivery, routingKey string, eventName string)

func connectToRabbitMQ(rmq *config.RabbitMQ) (*amqp.Connection, error) {
	// Connect to MQ
	log.Println("Connecting to MQ")
	conn, err := amqp.Dial(rmq.Url)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to MQ!")

	return conn, nil
}

func establishConnection(rmq *config.RabbitMQ, queueName string) (*amqp.Channel, <-chan amqp.Delivery, error) {

	// Create a new connection
	conn, err := connectToRabbitMQ(rmq)
	if err != nil {
		return nil, nil, err
	}

	// Create a new channel for this queue
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	msgs, err := ch.Consume(
		queueName,
		"processEngine", // consumer
		false,           // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		ch.Close()
		return nil, nil, fmt.Errorf("failed to register a consumer for queue %s: %v", queueName, err)
	}

	log.Println("Consuming from queue: " + queueName)

	return ch, msgs, nil
}

func Consume(rmq *config.RabbitMQ, queueName string, callback Callback) {

	forever := make(chan bool)

	go func() {
		for {
			ch, msgs, err := establishConnection(rmq, queueName)
			if err != nil {
				log.Printf("Failed to establish connection: %s", err)
				time.Sleep(time.Second * 5) // Wait before trying to reconnect
				continue
			}

			notify := ch.NotifyClose(make(chan *amqp.Error))

		consumeLoop:
			for {
				select {
				case err := <-notify:
					if err != nil {
						log.Printf("Channel closed for queue %s: %s", queueName, err)
					}
					ch.Close()
					break consumeLoop
				case d, ok := <-msgs:
					ctx := context.Background()
					if !ok {
						// The msgs channel has been closed, handle this situation, don't attempt to ack.
						// log.Println("Message channel closed")
						break consumeLoop
					}

					// Process the message
					log.Printf("Received a message from queue %s: %s", queueName, d.Body)

					// Get Routing Key
					routingKey, eventName := util.DetermineNewRoutingKeyAndEventName(d.RoutingKey, queueName)
					if routingKey == "error" {
						err := d.Ack(false)
						util.FailOnError(err, "Failed to ack")
						continue
					}

					// Process message
					callback(ch, ctx, d, routingKey, eventName)

					// Acknowledge the message
					err := d.Ack(false)
					util.FailOnError(err, "Failed to ack")
				}

			}
		}
	}()

	<-forever
}

func Publish(ch *amqp.Channel, ctx context.Context, msg []byte, routingKey string) {
	// err := ch.PublishWithContext(
	// 	ctx,
	// 	"topic.router", // exchange
	// 	routingKey,     // routing key
	// 	true,           // mandatory
	// 	false,          // immediate
	// 	amqp.Publishing{
	// 		ContentType: "text/plain",
	// 		Body:        msg,
	// 	})
	// util.FailOnError(err, "Failed to publish a message")
	// log.Printf("Published a message with routing key %s", routingKey)

}
