package chatbot

import (
	"SlitheringJake/pkg/markovchain"
	"github.com/gempir/go-twitch-irc/v3"
	"log"
	"sync"
)

func NewChatBot(config Config) (*ChatBot, error) {
	bot := ChatBot{
		Config:   config,
		Client:   nil,
		chains:   map[string]*markovchain.MarkovChain{},
		commands: map[string]CommandCallback{},
		mutexes: map[string]*sync.Mutex{
			// have to bootstrap this one, everything else should use .NewMutex
			"map_mutexes": {},
		},
	}

	bot.NewMutex("map_chains")
	bot.NewMutex("map_commands")

	bot.Client = twitch.NewClient(config.Twitch.Username, config.Twitch.Oauth)
	bot.Client.Join(config.Twitch.Channels...)

	bot.Client.OnPrivateMessage(bot.privateMessageHandler)
	bot.Client.OnConnect(func() {
		log.Println("Connected to IRC")
	})

	return &bot, nil
}

func (bot *ChatBot) Run() error {
	return bot.Client.Connect()
}
