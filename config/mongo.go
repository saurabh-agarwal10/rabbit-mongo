package config

import (
	"log"
	"os"

	mongosdk "bitbucket.org/kaleyra/mongo-sdk/mongo"
)

var MongoClient *mongosdk.Client

func MongoDBConnection() error {
	if MongoClient != nil {
		return nil
	}
	uri := mongosdk.URI{
		Username: os.Getenv("MONGODB_USER"),
		Password: os.Getenv("MONGODB_PASS"),
		Host:     os.Getenv("MONGODB_HOST"),
		Port:     os.Getenv("MONGODB_PORT"),
		DB:       os.Getenv("MONGODB_DBNAME"),
	}

	MongoClient, err = mongosdk.NewClient(uri)
	if err != nil {
		return err
	}

	err = MongoClient.Ping()
	if err != nil {
		return err
	}

	log.Println("MongoDB Successfully Connected.")
	return nil
}
