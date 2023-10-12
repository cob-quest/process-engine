package main

import (
	"log"

	"gitlab.com/cs302-2023/g3-team8/project/process-engine/config"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/util"
)

func main()  {
	log.Println("Loadingg .env file")
	config.InitEnvironment()
	log.Println(".env loaded!")

	rmq := config.SetupMQ()
	defer rmq.Conn.Close()
	defer rmq.Ch.Close()

	msgs, err := rmq.Ch.Consume(
		"queue.trigger.fromService",    // queue
		"processEngine",    // consumer
		false,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	util.FailOnError(err, "Failed to consume messages from queue.trigger.fromService")

	var forever chan struct{}

	go func ()  {
		log.Println("Consuming")
		for msg := range msgs {
			log.Printf("Received msg: %s\n", msg.Body)
	
			err = msg.Ack(false)
			if err != nil {
				log.Printf("Failed to ack message: %s", err)
			}
		}
	} ()

	log.Printf(" [*] Waiting for messages")
	<-forever
}