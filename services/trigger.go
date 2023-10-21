package services

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ProcessTrigger(ch *amqp.Channel, ctx context.Context, msg []byte, routingKey string) {
	Publish(ch, ctx, msg, routingKey)
}