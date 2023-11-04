package services

import (
	"context"
	"log"

	"github.com/davecgh/go-spew/spew"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/collections"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/util"
)

func ProcessImageBuilder(ch *amqp.Channel, ctx context.Context, d amqp.Delivery, routingKey string, eventName string) {

	// Unmarshal message
	m := util.UnmarshalJson(d.Body)
	log.Print("Json Body:")
	spew.Dump(m)

	// Process map
	temp := util.MapJsonToProcessEngine(m)
	log.Print("After Mapping:")
	spew.Dump(temp)

	temp.Event = &eventName

	// Store into db
	collections.CreateProcessEngine(temp)
}
