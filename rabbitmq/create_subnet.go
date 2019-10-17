package rabbitmq

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"hcc/harp/types"
)

func CreateSubnet(subnet types.Subnet) error {
	queue, err := Channel.QueueDeclare(
		"create_subnet",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return errors.New("failed to declare a create queue")
	}

	body, _ := json.Marshal(subnet)
	err = Channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing {
			ContentType:     "text/plain",
			ContentEncoding: "utf-8",
			Body:            body,
		})

	return nil
}
