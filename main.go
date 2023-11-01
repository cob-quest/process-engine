package main

import (
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/config"
	"gitlab.com/cs302-2023/g3-team8/project/process-engine/services"
)

func main() {

	rmq := config.SetupMQ()
	defer rmq.Conn.Close()
	defer rmq.Ch.Close()

	go services.Consume(rmq, "queue.challenge.fromService", services.ProcessChallenge)
	go services.Consume(rmq, "queue.notification.fromService", services.ProcessNotification)
	go services.Consume(rmq, "queue.platform.fromService", services.ProcessPlatform)
	go services.Consume(rmq, "queue.imageBuilder.fromService", services.ProcessImageBuilder)

	select {}
}