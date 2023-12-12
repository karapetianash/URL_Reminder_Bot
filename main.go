package main

import (
	"flag"
	"log"
)

func main() {
	t := mustToken()

	//tgClient := telegram.New(t)

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
