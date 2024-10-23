package slitheringjake

import (
	"context"
	"errors"
	"github.com/gempir/go-twitch-irc/v3"
	"log"
	"strings"
)

const MinimumSentenceLength = 6
const ErrorMessage = "Unable to fulfill."

// AuthCheck is a simple authorization callback, mods and VIPs can always used commands, otherwise its sub only
func (jake *SlitheringJake) AuthCheck(_ string, user twitch.User, message twitch.PrivateMessage) bool {
	if jake.Bot.IsModerator(user.Badges) || jake.Bot.IsVIP(user.Badges) {
		return true
	}

	if !jake.Bot.IsSubscriber(user.Badges) {
		return false
	}

	if !jake.checkLastUse(user.Name) {
		return false
	}

	// Assume if it hits here, it will reach a command so update the last use
	jake.updateLastUse(message.User.Name)
	return true
}

// GenerateCommand will send out a message to chat generated from the chat log based chain
func (jake *SlitheringJake) GenerateCommand(ctx context.Context, message twitch.PrivateMessage) error {
	chain := jake.Bot.GetChain(log_chain)
	defer jake.Bot.PutChain(log_chain)

	for {
		sentence, _ := chain.Generate()
		wordCount := strings.Count(sentence, " ") + 1
		if wordCount >= MinimumSentenceLength {
			jake.Bot.Client.Say(message.Channel, sentence)
			log.Printf("[*] %s", sentence)
			return nil
		}

		select {
		case <-ctx.Done():
			jake.deleteLastUse(message.User.Name)
			return ctx.Err()
		default:
			continue
		}
	}
}

// QuoteCommand will send out a message to chat generated from the quote based chain
func (jake *SlitheringJake) QuoteCommand(_ context.Context, message twitch.PrivateMessage) error {
	chain := jake.Bot.GetChain(quote_chain)
	defer jake.Bot.PutChain(quote_chain)

	sentence, _ := chain.Generate()
	jake.Bot.Client.Say(message.Channel, sentence)
	log.Printf("[*] %s", sentence)
	return nil
}

// ContainsCommand will send out a message generated from the log chain that contains the required terms
func (jake *SlitheringJake) ContainsCommand(ctx context.Context, message twitch.PrivateMessage) error {
	chain := jake.Bot.GetChain(log_chain)
	defer jake.Bot.PutChain(log_chain)

	words := strings.Split(strings.ToLower(message.Message), " ")
	if len(words) < 2 {
		jake.deleteLastUse(message.User.Name)
		return errors.New("missing arguments")
	}

	required := words[1:]
	for i := 0; i < len(required); i++ {
		required[i] = chain.NormalizeToken(required[i])
	}

	for {
		// We only actually use the first required term
		sentence, score := chain.Contains(required[0])

		if score > 0.0 {
			wordCount := strings.Count(sentence, " ") + 1
			if wordCount >= MinimumSentenceLength {
				jake.Bot.Client.Reply(message.Channel, message.ID, sentence)
				log.Printf("[*] %s", sentence)
				return nil
			}
		}

		select {
		case <-ctx.Done():
			jake.deleteLastUse(message.User.Name)
			return ctx.Err()
		default:
			continue
		}
	}
}

// QuestionCommand will try to generate a question for the streamer
func (jake *SlitheringJake) QuestionCommand(ctx context.Context, message twitch.PrivateMessage) error {
	chain := jake.Bot.GetChain(log_chain)
	defer jake.Bot.PutChain(log_chain)

	prefix := chain.NormalizeToken(jake.Config.StreamerName)
	for {
		sentence, _ := chain.StartsWith(prefix)
		if sentence[len(sentence)-1] == '?' {
			jake.Bot.Client.Say(message.Channel, sentence)
			log.Printf("[*] %s", sentence)
			return nil
		}

		select {
		case <-ctx.Done():
			jake.Bot.Client.Reply(message.Channel, message.ID, ErrorMessage)
			jake.deleteLastUse(message.User.Name)
			return ctx.Err()
		default:
			continue
		}
	}
}
