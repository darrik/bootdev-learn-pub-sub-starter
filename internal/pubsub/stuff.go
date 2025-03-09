package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	DurableQueueType SimpleQueueType = iota
	TransientQueueType
)

// simpleQueueType is an enum to represent "durable" or "transient"
func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, simpleQueueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	mqchan, err := conn.Channel()
	if err != nil {
		fmt.Printf("error creating rabbitmq channel: %s\n", err)
		return nil, amqp.Queue{}, err
	}

	durable := false
	autoDelete := false
	exclusive := false
	noWait := false

	switch simpleQueueType {
	case DurableQueueType:
		durable = true
	case TransientQueueType:
		autoDelete = true
		exclusive = true
	}

	queue, err := mqchan.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, nil)
	if err != nil {
		return nil, queue, err
	}

	err = mqchan.QueueBind(queueName, key, exchange, noWait, nil)
	if err != nil {
		return nil, queue, err
	}

	return mqchan, queue, nil
}
