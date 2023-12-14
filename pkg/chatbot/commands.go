package chatbot

import (
	"context"
	"github.com/gempir/go-twitch-irc/v3"
	"log"
	"strings"
	"time"
)

func (bot *ChatBot) AddCommand(command string, callback CommandCallback) {
	bot.Lock("map_commands")
	defer bot.Unlock("map_commands")

	command = strings.ToLower(command)
	bot.commands[command] = callback
}

func (bot *ChatBot) PrivateMessageHandler(message twitch.PrivateMessage) {
	go bot.dispatchCommand(message)
}

func (bot *ChatBot) dispatchCommand(message twitch.PrivateMessage) {
	if !strings.HasPrefix(message.Message, bot.Config.CommandPrefix) {
		return
	}
	words := strings.SplitN(message.Message, " ", 2)
	command := strings.ToLower(words[0][len(bot.Config.CommandPrefix):])

	if handler, found := bot.commands[command]; found {
		if bot.Config.AuthCheck != nil && !bot.Config.AuthCheck(command, message.User, message) {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() { cancel() }()

		if err := handler(ctx, message); err != nil {
			log.Printf("[%s] error: %s", command, err.Error())
		}
	}
}
