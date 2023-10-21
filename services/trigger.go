package services

import (
	"context"
	"log"

	"github.com/davecgh/go-spew/spew"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/collections"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/util"
)

func ProcessTrigger(ch *amqp.Channel, ctx context.Context, msg []byte, routingKey string) {

	// Unmarshal message
	temp := util.UnmarshalJson(msg)
	log.Print("Json Body:")
	spew.Dump(temp)

	// Store into db
	collections.CreateProcessEngine(temp)

	// Publish Message to Image Builder Topic
	Publish(ch, ctx, msg, routingKey)
}
