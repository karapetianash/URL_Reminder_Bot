package main

import (
	"flag"
	"log"

	tgClient "URLReminderBot/clients/telegram"
	"URLReminderBot/consumer/eventConsumer"
	"URLReminderBot/events/telegram"
	"URLReminderBot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	eventProcessor := telegram.New(
		tgClient.New(mustToken(), tgBotHost),
		files.New(storagePath))

	log.Println("service started")

	consumer := eventConsumer.New(eventProcessor, eventProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped")
	}
}

func mustToken() string {
	token := flag.String(
		"token-bot-token",
		"",
		"token for access to telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
