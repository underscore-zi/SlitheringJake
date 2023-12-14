package slitheringjake

import (
	"context"
	"github.com/gempir/go-twitch-irc/v3"
	"log"
	"strings"
	"time"
)

// NewMessage is a callback for new twitch messages. Command detection should happen in hte chatbot but we can add
// to the markov chain here.
func (jake *SlitheringJake) NewMessage(message twitch.PrivateMessage) {
	jake.Bot.PrivateMessageHandler(message)

	if jake.IsIgnoredAccount(message.User) {
		return
	}

	if jake.IsIgnoredPrefix(message.Message) {
		return
	}

	if jake.Config.MinimumMessageLength > len(message.Message) {
		return
	}

	logChain := jake.Bot.GetChain(log_chain)
	logChain.InsertLine(message.Message)
	jake.Bot.PutChain(log_chain)

	jake.Bot.Lock(log_mutex)
	if _, err := jake.logFile.WriteString(message.Message + "\n"); err != nil {
		log.Printf("[!] Failed to write to log file: %s", err)
	}
	jake.Bot.Unlock(log_mutex)

	jake.messageCount++
	if jake.messageCount >= jake.Config.MessageInterval {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := jake.GenerateCommand(ctx, message); err == nil {
			jake.messageCount = 0
		}
	}

}

// IsIgnoredAccount checks if the user is in the ignore list
func (jake *SlitheringJake) IsIgnoredAccount(user twitch.User) bool {
	// Ignore my own messages regardless of settings
	if user.Name == jake.Config.Twitch.Username {
		return true
	}

	for _, name := range jake.Config.Ignore.Accounts {
		if user.Name == name {
			return true
		}
	}
	return false
}

// IsIgnoredPrefix checks if the message starts with a prefix in the ignore list
func (jake *SlitheringJake) IsIgnoredPrefix(message string) bool {
	// Ignore own prefix regardless of settings
	if strings.HasPrefix(message, jake.Config.CommandPrefix) {
		return true
	}

	for _, prefix := range jake.Config.Ignore.Prefixes {
		if strings.HasPrefix(message, prefix) {
			return true
		}
	}

	return false
}
