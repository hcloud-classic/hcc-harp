package rabbitmq

import (
	"encoding/json"
	"hcc/harp/lib/dhcpd"
	"hcc/harp/lib/logger"
	"hcc/harp/model"
	"log"
)

// UpdateSubnet : Consume 'update_subnet' queues from RabbitMQ channel
func UpdateSubnet() error {
	qCreate, err := Channel.QueueDeclare(
		"update_subnet",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("update_subnet: Failed to declare a create queue")
		return err
	}

	msgsCreate, err := Channel.Consume(
		qCreate.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Logger.Println("update_subnet: Failed to register consumer")
		return err
	}

	go func() {
		for d := range msgsCreate {
			log.Printf("update_subnet: Received a create message: %s", d.Body)

			var subnet model.Subnet
			err = json.Unmarshal(d.Body, &subnet)
			if err != nil {
				logger.Logger.Println("update_subnet: Failed to unmarshal subnet data")
				return
			}

			//TODO: publish finished_update_subnet
		}
	}()

	return nil
}

// CreateDHCPDConfig : Consume 'create_dhcpd_config' queues from RabbitMQ channel
func CreateDHCPDConfig() error {
	qCreate, err := Channel.QueueDeclare(
		"create_dhcpd_config",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("create_dhcpd_config: Failed to declare a create queue")
		return err
	}

	msgsCreate, err := Channel.Consume(
		qCreate.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Logger.Println("create_dhcpd_config: Failed to register consumer")
		return err
	}

	go func() {
		for d := range msgsCreate {
			log.Printf("create_dhcpd_config: Received a create message: %s", d.Body)

			var subnet model.Subnet
			err = json.Unmarshal(d.Body, &subnet)
			if err != nil {
				logger.Logger.Println("create_dhcpd_config: Failed to unmarshal subnet data")
				return
			}

			err := dhcpd.CreateConfig(subnet.UUID, subnet.NetworkIP, subnet.Netmask, subnet.Gateway,
				subnet.NextServer, subnet.NameServer, subnet.DomainName,
				subnet.MaxNodes, subnet.NodeUUIDs, subnet.LeaderNodeUUID, subnet.OS, subnet.Name)
			if err != nil {
				logger.Logger.Println("create_dhcpd_config: " + err.Error())
				return
			}

			//logger.Logger.Println("create_dhcpd_config: UUID = " + uuid + ": " + result)
		}
	}()

	return nil
}
