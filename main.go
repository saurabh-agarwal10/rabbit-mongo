package main

import (
	"encoding/json"
	"log"
	"time"

	"rabbitmongo-server/config"
	"rabbitmongo-server/producing"

	"github.com/joho/godotenv"
)

//declaring struct for body of queue
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

	//user details to send to rmq queue
	jsonBody := &Body{
		Name:          "Saurabh",
		Address:       "Bengaluru",
		Subscriptions: "Calls API",
		Timestamp:     time.Now().Unix(),
	}
	//parsing user details to json
	userDetails, err := json.Marshal(jsonBody)
	if err != nil {
		log.Fatal("Couldn't marshal struct to json")
		return
	}

	//publishing to rmq
	err = producing.PublishToQueue("user_details_queue", userDetails)
	if err != nil {
		log.Fatal("Couldn't publish to queue")
		return
	}

}
