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

	gs := gamelogic.NewGameState(username)
	for {
		input := gamelogic.GetInput()
		if len(input) < 1 || len(input[0]) < 1 {
			continue
		}

		if input[0] == "spawn" {
			err := gs.CommandSpawn(input)
			if err != nil {
				fmt.Printf("error: %s\n", err)
			}
		} else if input[0] == "move" {
			_, err := gs.CommandMove(input)
			if err != nil {
				fmt.Printf("error: %s\n", err)
			}
		} else if input[0] == "status" {
			gs.CommandStatus()
		} else if input[0] == "help" {
			gamelogic.PrintClientHelp()
		} else if input[0] == "spam" {
			fmt.Println("Spamming not allowed yet!")
		} else if input[0] == "quit" {
			gamelogic.PrintQuit()
			return
		} else {
			fmt.Println("Unknown command")
		}
	}

	// fmt.Println("Press CTRL-C to quit.")
	// // wait for ctrl+c
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
}
