package chatbot

import (
	"context"
	"github.com/gempir/go-twitch-irc/v3"
	"log"
	"strings"
)

func (bot *ChatBot) generateHandler(ctx context.Context, message twitch.PrivateMessage) error {
	chain := bot.GetChain("log")
	defer bot.PutChain("log")

	if !bot.IsSubscriber(message.User.Badges) {
		return nil
	}

	for {
		sentence := chain.GenerateSentence()
		if strings.Count(sentence, " ") >= 5 {
			bot.Client.Say(message.Channel, sentence)
			log.Printf("[*] %s", sentence)
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			continue
		}
	}
}

func (bot *ChatBot) containsHandler(ctx context.Context, message twitch.PrivateMessage) error {
	chain := bot.GetChain("log")
	defer bot.PutChain("log")

	if !bot.IsSubscriber(message.User.Badges) {
		return nil
	}

	words := strings.Split(strings.ToLower(message.Message), " ")
	if len(words) < 2 {
		return nil
	}
	required := words[1:]

	for {
		sentence := chain.GenerateSentence()
		lower := strings.ToLower(sentence)

		found := true
		for _, word := range required {
			// Try and prevent matching where a word appearsa in the middle of another word,
			// but still catch minor variations like plural or verb endings
			if !strings.HasPrefix(lower, word) && !strings.Contains(lower, " "+word) {
				found = false
				break
			}
		}

		if found {
			bot.Client.Reply(message.Channel, message.ID, sentence)
			log.Printf("[*] %s", sentence)
			return nil
		}

		select {
		case <-ctx.Done():
			bot.Client.Reply(message.Channel, message.ID, "Unable to fulfill")
			return ctx.Err()
		default:
			continue
		}
	}
}

func (bot *ChatBot) questionHandler(ctx context.Context, message twitch.PrivateMessage) error {
	chain := bot.GetChain("log")
	defer bot.PutChain("log")

	if !bot.IsSubscriber(message.User.Badges) {
		return nil
	}

	streamer := strings.ToLower(bot.Config.StreamerName)
	for {
		sentence := chain.GenerateSentence()
		lower := strings.ToLower(sentence)

		if lower[len(lower)-1] == '?' && (strings.HasPrefix(lower, streamer) || strings.Contains(lower, " "+streamer)) {
			bot.Client.Say(message.Channel, sentence)
			log.Printf("[*] %s", sentence)
			return nil
		}

		select {
		case <-ctx.Done():
			bot.Client.Reply(message.Channel, message.ID, "Unable to fulfill")
			return ctx.Err()
		default:
			continue
		}
	}
}
