package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"hcc/harp/lib/logger"
	"hcc/harp/model"
)

// GetNodes : Publish 'get_nodes' queues to RabbitMQ channel
func GetNodes(nodeNr int, serverUUID string) error {
	qCreate, err := Channel.QueueDeclare(
		"get_nodes",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("get_nodes: Failed to declare a create queue")
		return err
	}

	var server model.Server
	server.UUID = serverUUID
	server.NodeNr = nodeNr

	body, _ := json.Marshal(server)
	err = Channel.Publish(
		"",
		qCreate.Name,
		false,
		false,
		amqp.Publishing {
			ContentType:     "text/plain",
			ContentEncoding: "utf-8",
			Body:            body,
		})
	if err != nil {
		logger.Logger.Println("get_nodes: Failed to register publisher")
		return err
	}

	return nil
}
