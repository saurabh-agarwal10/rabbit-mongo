package main

import (
	"encoding/json"
	"log"

	"rabbitmongo-server/config"

	"github.com/joho/godotenv"
)

type Body struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	Subscriptions string `json:"subscriptions"`
	Timestamp     int64  `json:"timestamp"`
}

func main() {
	// Load the environment file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("unable to load .env file ")
	}

	//connecting to rmq
	err = config.InitRabbitMQ()
	if err != nil {
		log.Fatal("unable to initialize rabbitmq", err)
	}

	// disconnect rmq connection and channel on exit
	defer config.RMQDisconnect()

	// mongo connection
	err = config.MongoDBConnection()
	if err != nil {
		log.Fatal("unable to connect mongodb ")
	}

	//close mongo connection
	defer func() {
		_ = config.MongoClient.Disconnect()
	}()

	//Consume data from rmq
	userDetails, err := config.RMQChan.Consume(
		"user_details_queue", // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		log.Fatal("Couldn't consume data from queue.")
		return
	}

	for m := range userDetails {

		detail := &Body{}

		err := json.Unmarshal(m.Body, detail)
		if err != nil {
			log.Println("Couldn't Unmarshal data to json")
		}

		collection := config.MongoClient.Collection("user_details")

		result, err := collection.InsertOne(detail)
		if err != nil {
			log.Println("Couldn't publish data to MongoDB")
		}

		log.Printf("Inserted document with _id: %v\n", result)

	}
}
