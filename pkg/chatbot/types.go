package chatbot

import (
	"SlitheringJake/pkg/mctg"
	"context"
	"github.com/gempir/go-twitch-irc/v3"
	"sync"
)

type ChatBot struct {
	Config       Config
	Client       *twitch.Client
	chains       map[string]*mctg.MCTG
	commands     map[string]CommandCallback
	mutexes      map[string]*sync.Mutex
	messageCount int
}

type Config struct {
	// LogFile is the filename that the logged messages will be written to and read from for the markov chain
	LogFile string
	// CommandPrefix is the string that will prefix all commands. This must be unique from any in Ignore.Prefixes
	CommandPrefix string

	// Twitch connection settings for twitch
	Twitch struct {
		// Username the username to authenticate with on the twitch IRC
		Username string
		// Oauth is the oauth token used to authenticate with twitch IRC
		Oauth string
		// Channels is a list of channels to join
		Channels []string
	}

	Ignore struct {
		// Accounts should contains all accounts thaat should be ignored like nightbot or co2_bot
		Accounts []string
		// Prefixes are any tokens used by _other_ bots to prefix their commands, like `!` and `.`, these can be multicharacter
		Prefixes []string
	}

	// MinimumMessageLength for a message to be logged
	MinimumMessageLength int
	// MessageInterval is the frequency for which messages will be automatically generated
	MessageInterval int
	// StreamerName is used by the -question command to find messages asking the streamer a question
	StreamerName string
}

type CommandCallback func(context context.Context, message twitch.PrivateMessage) error
