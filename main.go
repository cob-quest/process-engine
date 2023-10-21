package main

import (
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/config"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/services"
)

func main() {

	rmq := config.SetupMQ()
	defer rmq.Conn.Close()
	defer rmq.Ch.Close()

	// go services.Consume(rmq, "queue.assignment.fromService")
	// go services.Consume(rmq, "queue.notification.fromService")
	go services.Consume(rmq, "queue.trigger.fromService", services.ProcessTrigger)
	go services.Consume(rmq, "queue.imageBuilder.fromService", services.ProcessImageBuilder)

	select {}
}