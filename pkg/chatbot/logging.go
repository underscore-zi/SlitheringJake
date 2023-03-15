package chatbot

import (
	"context"
	"fmt"
	"github.com/gempir/go-twitch-irc/v3"
	"log"
	"os"
	"strings"
	"time"
)

func (bot *ChatBot) isIgnoredAccount(user twitch.User) bool {
	for _, name := range bot.Config.Ignore.Accounts {
		if user.Name == name {
			return true
		}
	}
	return false
}

func (bot *ChatBot) isIgnoredPrefix(message twitch.PrivateMessage) bool {
	for _, prefix := range bot.Config.Ignore.Prefixes {
		if strings.HasPrefix(message.Message, prefix) {
			return true
		}
	}
	return false
}

func (bot *ChatBot) isPotentialCommand(message twitch.PrivateMessage) bool {
	return strings.HasPrefix(message.Message, bot.Config.CommandPrefix)
}

func (bot *ChatBot) logMessage(message twitch.PrivateMessage) {
	if bot.isIgnoredAccount(message.User) || bot.isIgnoredPrefix(message) {
		return
	}
	if bot.isPotentialCommand(message) {
		return
	}
	if bot.Config.MinimumMessageLength > len(message.Message) {
		return
	}

	bot.lock("file_log")
	defer bot.unlock("file_log")

	fp, err := os.OpenFile(bot.Config.LogFile, os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	defer func() { _ = fp.Close() }()

	_, err = fp.WriteString(message.Message + "\n")
	if err != nil {
		log.Printf("[!] Failed to write to log file: %s", err)
	}

	chain := bot.GetChain("log")
	defer bot.PutChain("log")
	chain.InsertLine(message.Message)
}

func (bot *ChatBot) privateMessageHandler(message twitch.PrivateMessage) {
	log.Println(message.Message)
	go bot.logMessage(message)
	go bot.dispatchCommand(message)

	if bot.isIgnoredPrefix(message) || bot.isIgnoredAccount(message.User) || bot.isPotentialCommand(message) {
		return
	}
	bot.messageCount++
	if bot.messageCount >= bot.Config.MessageInterval {
		bot.messageCount = 0

		// Reuse the -generate command by spoofing it here
		message.User.Badges["broadcaster"] = 1
		message.Message = fmt.Sprintf("%sgenerate", bot.Config.CommandPrefix)
		go bot.dispatchCommand(message)
	}
}

func (bot *ChatBot) dispatchCommand(message twitch.PrivateMessage) {
	if !bot.isPotentialCommand(message) {
		return
	}
	words := strings.SplitN(message.Message, " ", 2)
	command := strings.ToLower(words[0][len(bot.Config.CommandPrefix):])

	if handler, found := bot.commands[command]; found {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() { cancel() }()

		if err := handler(ctx, message); err != nil {
			log.Printf("[%s] error: %s", command, err.Error())
		}
	}
}
