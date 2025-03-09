package main

import (
	"fmt"
	// "os"
	// "os/signal"

	"github.com/darrik/bootdev-learn-pub-sub-starter/internal/gamelogic"
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

	_, _, err = pubsub.DeclareAndBind(dial, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.DurableQueueType)
	if err != nil {
		fmt.Printf("rabbitmq error: %s\n", err)
		return
	}

	mqchan, err := dial.Channel()
	if err != nil {
		fmt.Printf("error creating rabbitmq channel: %s\n", err)
		return
	}

	// fmt.Println("Press CTRL-C to quit.")

	gamelogic.PrintServerHelp()
	for {
		input := gamelogic.GetInput()
		if len(input) < 1 || len(input[0]) < 1 {
			continue
		}

		if input[0] == "pause" {
			fmt.Println("Pausing game")

			err = pubsub.PublishJSON(mqchan, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{IsPaused: true})
			if err != nil {
				fmt.Printf("error pausing game: %s\n", err)
				return
			}
		} else if input[0] == "resume" {
			fmt.Println("Resuming game")

			err = pubsub.PublishJSON(mqchan, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{IsPaused: false})
			if err != nil {
				fmt.Printf("error unpausing game: %s\n", err)
				return
			}
		} else if input[0] == "quit" {
			fmt.Println("Quitting...")

			return
		} else {
			fmt.Println("Unknown command")
		}
	}

	// // wait for ctrl+c
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
}
