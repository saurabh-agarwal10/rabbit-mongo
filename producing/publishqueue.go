package producing

import (
	"log"

	"rabbitmongo-server/config"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Channel *amqp.Channel
}

func PublishToQueue(qName string, rmqBody []byte) error {
	log.Printf("Publishing to %s queue", qName)

	// publish to rmq
	err := config.RMQChan.Publish(
		"USER_DETAILS_EXCHANGE",    // exchange
		qName, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        rmqBody,
			// Type:        publishType,
		})

	if err != nil {
		log.Fatal("Error while publishing",err)
		return err
	}

	return nil
}

