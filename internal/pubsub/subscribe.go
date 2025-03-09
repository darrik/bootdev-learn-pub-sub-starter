package pubsub

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](conn *amqp.Connection, exchange string, queueName string, key string, simpleQueueType SimpleQueueType, handler func(T)) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	deliveryChannel, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range deliveryChannel {
			var data T
			err := json.Unmarshal(delivery.Body, &data)
			if err != nil {
				panic(err)
			}
			handler(data)
			err = delivery.Ack(false)
			if err != nil {
				panic(err)
			}
		}
	}()

	return nil
}
