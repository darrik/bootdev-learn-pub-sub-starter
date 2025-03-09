package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/darrik/bootdev-learn-pub-sub-starter/internal/gamelogic"
	"github.com/darrik/bootdev-learn-pub-sub-starter/internal/pubsub"
	"github.com/darrik/bootdev-learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

const connStr = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril client...")

	dial, err := amqp.Dial(connStr)
	if err != nil {
		fmt.Printf("error connecting to rabbitmq: %s\n", err)
		return
	}
	defer dial.Close()
	fmt.Println("Connection to RabbitMQ established!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	_, _, err = pubsub.DeclareAndBind(dial, routing.ExchangePerilDirect, routing.PauseKey+"."+username, routing.PauseKey, pubsub.TransientQueueType)
	if err != nil {
		fmt.Printf("rabbitmq error: %s\n", err)
		return
	}

	fmt.Println("Press CTRL-C to quit.")
	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
