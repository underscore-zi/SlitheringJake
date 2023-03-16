package chatbot

import (
	"context"
	"github.com/gempir/go-twitch-irc/v3"
	"log"
	"strings"
	"time"
)

func (bot *ChatBot) generateHandler(ctx context.Context, message twitch.PrivateMessage) error {
	chain := bot.GetChain("log")
	defer bot.PutChain("log")

	if !bot.IsSubscriber(message.User.Badges) {
		return nil
	}
	if !bot.CheckLastUse(message.User, time.Minute*2) {
		return nil
	}

	for {
		sentence, _ := chain.Generate()
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

func (bot *ChatBot) resetLastUser(user twitch.User) {
	delete(bot.lastUse, user.Name)
}

func (bot *ChatBot) CheckLastUse(user twitch.User, duration time.Duration) bool {
	if bot.IsModerator(user.Badges) || bot.IsVIP(user.Badges) {
		return true
	}

	if last, found := bot.lastUse[user.Name]; found {
		nextUse := last.Add(duration)
		if !time.Now().After(nextUse) {
			return false
		}
	}

	bot.lastUse[user.Name] = time.Now()
	return true
}

func (bot *ChatBot) containsHandler(ctx context.Context, message twitch.PrivateMessage) error {
	chain := bot.GetChain("log")
	defer bot.PutChain("log")

	if !bot.IsSubscriber(message.User.Badges) {
		return nil
	}
	if !bot.CheckLastUse(message.User, time.Minute*2) {
		return nil
	}

	words := strings.Split(strings.ToLower(message.Message), " ")
	if len(words) < 2 {
		return nil
	}
	required := words[1:]

	for {
		sentence, _ := chain.Generate()
		lower := strings.ToLower(sentence)

		found := true
		for _, word := range required {
			// Try and prevent matching where a word appears in the middle of another word,
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
			bot.resetLastUser(message.User)
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
	if !bot.CheckLastUse(message.User, time.Minute*2) {
		return nil
	}

	streamer := strings.ToLower(bot.Config.StreamerName)
	prefix := chain.NormalizeToken(streamer)

	for {
		sentence, _ := chain.StartsWith(prefix)
		if sentence[len(sentence)-1] == '?' {
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
