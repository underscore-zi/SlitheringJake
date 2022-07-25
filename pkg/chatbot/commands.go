package chatbot

import "strings"

func (bot *ChatBot) AddCommand(command string, callback CommandCallback) {
	bot.lock("map_commands")
	defer bot.unlock("map_commands")

	command = strings.ToLower(command)
	bot.commands[command] = callback
}
