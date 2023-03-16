package main

import (
	"context"
	"github.com/gempir/go-twitch-irc/v3"
	"log"
)

func QuoteCommand(_ context.Context, message twitch.PrivateMessage) error {
	chain := bot.GetChain("quotes")
	defer bot.PutChain("quotes")

	if !bot.IsSubscriber(message.User.Badges) {
		return nil
	}

	sentence, _ := chain.Generate()
	bot.Client.Say(message.Channel, sentence)
	log.Printf("[*] %s", sentence)
	return nil
}
