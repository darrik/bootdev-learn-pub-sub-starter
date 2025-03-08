package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange string, key string, val T) error {
  jsonval, err := json.Marshal(val)
  if err != nil {
    return err
  }

  err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{ContentType: "application/json", Body: jsonval})
  if err != nil {
    return err
  }

  return nil
}
