package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/modeckrus/scheduler"
	"github.com/modeckrus/scheduler/storage"
)

func TaskWithoutArgs() {
	log.Println("TaskWithoutArgs is executed")
}

func TaskWithArgs(message string) {
	log.Println("TaskWithArgs is executed. message:", message)
}

type Payload struct {
	Name string `json:"name"`
}

func NotificationStorage(lol string) {
	payload := Payload{}
	err := json.Unmarshal([]byte(lol), &payload)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(payload)
}

func main() {
	storage := storage.NewMongoDBStorage(
		storage.MongoDBConfig{
			ConnectionUrl: "mongodb://localhost:27017/?readPreference=primary&ssl=false&connect=direct",
			Db:            "example",
		},
	)

	if err := storage.Connect(); err != nil {
		log.Fatal("Could not connect to db", err)
	}

	if err := storage.Initialize(); err != nil {
		log.Fatal("Could not intialize database", err)
	}
	s := scheduler.New(storage)
	fmt.Println("Storage Initialized")
	// Start a task without arguments
	if _, err := s.RunAfter(2*time.Second, TaskWithoutArgs); err != nil {
		log.Fatal(err)
	}
	payload := Payload{
		Name: "name",
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	// Start a task with arguments
	if _, err := s.RunEvery(5*time.Second, NotificationStorage, string(bytes)); err != nil {
		log.Fatal(err)
	}
	s.Start()
	s.Wait()
}
