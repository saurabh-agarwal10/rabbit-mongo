package config

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

var RMQConn *amqp.Connection
var RMQChan *amqp.Channel
var err error

func rabbitMQConnection() (*amqp.Connection, error) {
	Host := os.Getenv("RABBITMQ_HOST")
	Port := os.Getenv("RABBITMQ_PORT")
	Username := os.Getenv("RABBITMQ_USER")
	Password := os.Getenv("RABBITMQ_PASS")

	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s", Username, Password, Host, Port)

	log.Println(connectionString)
	conn, err := amqp.Dial(connectionString)

	if err != nil {
		log.Println("RabbitMQ Connection failed")
		return nil, err
	}

	log.Println("RabbitMQ Successfully Connected.")

	return conn, nil
}

func RMQDisconnect() {
	_ = RMQConn.Close()
	_ = RMQChan.Close()

	log.Println("RMQ connections and Channels closed.")
}

func InitRabbitMQ() error {

	if RMQConn == nil || RMQConn.IsClosed() {
		RMQConn, err = rabbitMQConnection()
		if err != nil {
			return err
		}
	}

	// rabbitmq channel initialization
	RMQChan, err = RMQConn.Channel()
	if err != nil {
		return err
	}

	log.Println("new RMQ connection and channel successfully created.")

	// now declaring all the required queues and exchanges
	return declareQueuesAndExchanges()
}

// function to declare all the queues, exchanges and bindings between them
func declareQueuesAndExchanges() error {
	err := declareQueues()
	if err != nil {
		return err
	}
	err = declareExchanges()
	if err != nil {
		return err
	}

	err = bindQueueWithExchange()
	if err != nil {
		return err
	}

	log.Println("All queues and exchanges has been declared successfully")
	return nil
}

// function to declare queues
func declareQueues() error {
	//if queue declare properties are same as used in this function.
	queueList := []string{
		"user_details_queue",
	}

	for _, queue := range queueList {
		_, err := RMQChan.QueueDeclare(
			queue, //name
			true,  //durable
			false, //auto-delete
			false, //exclusive
			false, //noWait
			nil,   //args
		)

		if err != nil {
			log.Fatal("Error declaring queue", err)
			return err
		}
	}
	return nil
}

// function to declare exchanges
func declareExchanges() error {
	//if exchange declare properties are same as used in this function. */
	queueExchangeList := []string{"USER_DETAILS_EXCHANGE"}

	for _, exchange := range queueExchangeList {
		err := RMQChan.ExchangeDeclare(
			exchange, // name
			"fanout", // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)

		if err != nil {
			log.Fatal("Failed to declare exchange.", err)
			return err
		}
	}
	return nil
}

// function to bind the exchanges with queue
func bindQueueWithExchange() error {
	//Note: All the binding between queue and exchange should come here
	queuewithExchange := map[string]string{"user_details_queue": "USER_DETAILS_EXCHANGE"}

	for queue, exchange := range queuewithExchange {
		err = RMQChan.QueueBind(
			queue,    // queue name
			"",       // routing key
			exchange, // exchange
			false,
			nil,
		)

		if err != nil {
			log.Fatal("Error occurred while binding queue with exchange", err)
			return err
		}
	}
	return nil
}

func IsRMQConnClosed() bool {
	return (RMQConn != nil && RMQConn.IsClosed())
}
