package main

import (
	"log"

	"gitlab.com/cs302-2023/g3-team8/project/process-engine/config"
)

func main()  {
	log.Println("Loadingg .env file")
	config.InitEnvironment()
	log.Println(".env loaded!")

	rmq := config.SetupMQ("testQueue")
	defer rmq.Conn.Close()
	defer rmq.Channel.Close()
}