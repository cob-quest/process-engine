package services

import (
	"context"
	"log"

	"github.com/davecgh/go-spew/spew"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/collections"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/util"
)

func ProcessImageBuilder(ch *amqp.Channel, ctx context.Context, msg []byte, routingKey string) {

	// Unmarshal message
	temp := util.UnmarshalJson(msg)
	log.Print("Json Body:")
	spew.Dump(temp)

	// Store into db
	collections.CreateProcessEngine(temp)
}
