package chatbot

import (
	"SlitheringJake/pkg/markovchain"
	"context"
	"github.com/gempir/go-twitch-irc/v3"
	"sync"
)

type ChatBot struct {
	Config   Config
	Client   *twitch.Client
	chains   map[string]*markovchain.MarkovChain
	commands map[string]CommandCallback
	mutexes  map[string]*sync.Mutex
}

type TwitchConfig struct {
	// Username the username to authenticate with on the twitch IRC
	Username string `json:"username"`

	// Oauth is the oauth token used to authenticate with twitch IRC
	Oauth string `json:"oauth"`

	// Channels is a list of channels to join
	Channels []string `json:"channels"`
}

type Config struct {
	// CommandPrefix is the string that will prefix all commands. This must be unique from any in Ignore.Prefixes
	CommandPrefix string
	// Twitch connection settings for twitch
	Twitch TwitchConfig
	// AuthCheck is run before a command is dispatched to check if the user is authorized
	AuthCheck AuthCallback
}

type CommandCallback func(context context.Context, message twitch.PrivateMessage) error
type AuthCallback func(command string, user twitch.User, message twitch.PrivateMessage) bool
