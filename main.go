package main

import (
	"flag"
	"log"

	"URLReminderBot/clients/telegram"
)

const tgBotHost = "api.telegram.org"

func main() {
	tgClient := telegram.New(mustToken(), tgBotHost)

	//fetcher := fetcher.New(tgClient)

	//processor := processor.New(tgClient)

	//consumer := consumer.Start(fetcher, processor)
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
