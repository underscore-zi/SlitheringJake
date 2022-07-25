package chatbot

import (
	"SlitheringJake/pkg/mctg"
	"github.com/gempir/go-twitch-irc/v3"
	"log"
	"sync"
)

func NewChatBot(config Config) (*ChatBot, error) {
	bot := ChatBot{
		Config:   config,
		Client:   nil,
		chains:   map[string]*mctg.MCTG{},
		commands: map[string]CommandCallback{},
		mutexes: map[string]*sync.Mutex{
			// have to bootstrap this one, everything else should use .newMutex
			"map_mutexes": &sync.Mutex{},
		},
	}

	bot.newMutex("map_chains")
	bot.newMutex("map_commands")
	bot.newMutex("file_log")

	bot.NewChain("log", config.LogFile, 2)

	bot.Client = twitch.NewClient(config.Twitch.Username, config.Twitch.Oauth)
	bot.Client.Join(config.Twitch.Channels...)

	bot.Client.OnPrivateMessage(bot.privateMessageHandler)
	bot.Client.OnConnect(func() {
		log.Println("Connected to IRC")
	})

	bot.AddCommand("generate", bot.generateHandler)
	bot.AddCommand("contains", bot.containsHandler)
	bot.AddCommand("question", bot.questionHandler)

	return &bot, nil
}

func (bot *ChatBot) Run() error {
	return bot.Client.Connect()
}
