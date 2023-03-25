package slitheringjake

import (
	"SlitheringJake/pkg/chatbot"
	"os"
	"time"
)

const lastuse_mutex = "lastuse_mutex"
const log_mutex = "log_mutex"
const log_chain = "log"
const quote_chain = "quote"

type Config struct {
	// LogFile is the path to the log file to be used by the bot
	LogFile string `json:"log_file"`

	// QuoteFile contains the new-line delimited quotes to be used by the -quote command
	QuoteFile string `json:"quote_file"`

	// MinimumMessageLength for a message to be logged
	MinimumMessageLength int `json:"minimum_message_length"`

	// MessageInterval is the frequency for which messages will be automatically generated
	MessageInterval int `json:"message_interval"`

	// StreamerName should reflect what people actually call the streamer in chat, not just their username
	StreamerName string `json:"streamer_name"`

	Ignore struct {
		// Accounts should contains all accounts that should be ignored like nightbot or co2_bot
		Accounts []string `json:"accounts"`
		// Prefixes are any tokens used by _other_ bots to prefix their commands, like `!` and `.`, these can be multicharacter
		Prefixes []string `json:"prefixes"`
	} `json:"ignore"`

	// Twitch is auth information for the bot's twitch chat account
	Twitch chatbot.TwitchConfig `json:"twitch"`

	// UseInterval is the frequency for which subs can use commands
	UseInterval int `json:"use_interval"`

	// CommandPrefix is the token used to indicate the start of a command
	CommandPrefix string `json:"command_prefix"`
}

type SlitheringJake struct {
	Bot          *chatbot.ChatBot
	Config       Config
	logFile      *os.File
	lastUse      map[string]time.Time
	messageCount int
}

// NewSlitheringJake creates a new instance of the SlitheringJake bot
func NewSlitheringJake(config Config) (*SlitheringJake, error) {
	var jake SlitheringJake
	var err error

	jake.Config = config

	botConfig := chatbot.Config{
		CommandPrefix: jake.Config.CommandPrefix,
		Twitch:        config.Twitch,
		AuthCheck:     jake.AuthCheck,
	}

	if jake.Bot, err = chatbot.NewChatBot(botConfig); err != nil {
		return nil, err
	}

	jake.Bot.NewMutex(lastuse_mutex) // controls access to lastuse map
	jake.Bot.NewMutex(log_mutex)     // to control log file writes
	jake.Bot.NewChain(log_chain, config.LogFile, 2)

	if jake.logFile, err = os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return nil, err
	}

	jake.Bot.AddCommand("generate", jake.GenerateCommand)
	jake.Bot.AddCommand("question", jake.QuestionCommand)
	jake.Bot.AddCommand("contains", jake.ContainsCommand)

	if jake.Config.QuoteFile != "" {
		jake.Bot.NewChain(quote_chain, jake.Config.QuoteFile, 2)
		jake.Bot.AddCommand("quote", jake.QuoteCommand)
	}

	jake.Bot.Client.OnPrivateMessage(jake.NewMessage)

	return &jake, nil
}
