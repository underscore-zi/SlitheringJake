package main

import (
	"context"
	"github.com/gempir/go-twitch-irc/v3"
	"log"
	"strings"
)

func QuoteCommand(ctx context.Context, message twitch.PrivateMessage) error {
	chain := bot.GetChain("quotes")
	defer bot.PutChain("quotes")

	if !bot.IsSubscriber(message.User.Badges) {
		return nil
	}

	for {
		sentence, uniq := chain.Generate()
		if strings.Count(sentence, " ") >= 5 && uniq > 1.5 {
			found := false
			for _, quote := range allQuotes {
				if strings.Contains(quote, sentence) {
					found = true
					break
				} else {

				}
			}

			if !found {
				bot.Client.Say(message.Channel, sentence)
				log.Printf("[*] %s", sentence)
				return nil
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			continue
		}
	}
}
