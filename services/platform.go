package services

import (
	"context"
	"log"

	"github.com/davecgh/go-spew/spew"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/collections"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/util"
)

func ProcessPlatform(ch *amqp.Channel, ctx context.Context, msg []byte, routingKey string, eventName string) {

	// Unmarshal message
	m := util.UnmarshalJson(msg)
	log.Print("Json Body:")
	spew.Dump(m)

	// Process map
	temp := util.MapJsonToProcessEngine(m)
	log.Print("After Mapping:")
	spew.Dump(temp)

	temp.Event = &eventName

	// Store into db
	collections.CreateProcessEngine(temp)

	// Publish Message to Image Builder Topic
	Publish(ch, ctx, msg, routingKey)
}
