package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/darrik/bootdev-learn-pub-sub-starter/internal/pubsub"
	"github.com/darrik/bootdev-learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

const connStr = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril server...")
	dial, err := amqp.Dial(connStr)
	if err != nil {
		fmt.Printf("error connecting to rabbitmq: %s\n", err)
		return
	}
	defer dial.Close()
	fmt.Println("Connection to RabbitMQ established!")
	fmt.Println("Press CTRL-C to quit.")

	mqchan, err := dial.Channel()
	if err != nil {
		fmt.Printf("error creating rabbitmq channel: %s\n", err)
		return
	}

	err = pubsub.PublishJSON(mqchan, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{IsPaused: true})
	if err != nil {
		fmt.Printf("error publishing json: %s\n", err)
		return
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
