package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	tgClient "ReadAdvisor/clients/telegram"
	"ReadAdvisor/consumer/event_consumer"
	"ReadAdvisor/events/telegram"
	"ReadAdvisor/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "/home/olekdsii/Desktop/data/sqlite/storage.db"
	batchSize         = 100
)

type Configs struct {
	Token string `json:"token"`
}

func main() {
	// s:= files.New(storagePath)
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatalf("can`t connect to storage: %v", err)
	}
	if err = s.Init(context.TODO()); err != nil {
		log.Fatalf("can`t init storage: %v", err)

	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, getConfigs()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, 100)

	if err := consumer.Start(); err != nil {
		log.Fatal("service was stopped", err)
	}

}

func getConfigs() (token string) {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		panic("Config file was not provided")
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configs Configs
	json.Unmarshal(byteValue, &configs)

	if configs.Token == "" {
		panic("No token provided")

	}
	return configs.Token
}
