package main

import (
	"SlitheringJake/pkg/chatbot"
	"context"
	"encoding/json"
	"flag"
	"github.com/gempir/go-twitch-irc/v3"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var bot *chatbot.ChatBot

func main() {
	var err error
	var config chatbot.Config
	var configfn, quotesfn string
	flag.StringVar(&configfn, "config", "", "Path to config file (json)")
	flag.StringVar(&quotesfn, "quotes", "", "Path to quotes log file (optional)")
	flag.Parse()

	if configfn == "" {
		flag.Usage()
		return
	}

	if content, err := ioutil.ReadFile(configfn); err == nil {
		if err := json.Unmarshal(content, &config); err != nil {
			log.Printf("failed to parse config: %s", err.Error())
			return
		}
	} else {
		log.Printf("failed to load config: %s", err.Error())
		return
	}

	bot, err = chatbot.NewChatBot(config)
	if err != nil {
		log.Printf("Failed to create bot: %s", err.Error())
		return
	}

	if quotesfn != "" {
		if _, err := os.Stat(quotesfn); err == nil {
			bot.NewChain("quotes", quotesfn, 3)
			bot.AddCommand("quote", QuoteCommand)
		} else {
			log.Printf("Quotes file was provided, but does not exist.")
			return
		}
	}
	err = bot.Run()
	if err != nil {
		panic(err)
	}

}

func QuoteCommand(ctx context.Context, message twitch.PrivateMessage) error {
	chain := bot.GetChain("quotes")
	defer bot.PutChain("quotes")

	if !bot.IsSubscriber(message.User.Badges) {
		return nil
	}

	for {
		sentence := chain.GenerateSentence()
		if strings.Count(sentence, " ") >= 5 {
			bot.Client.Say(message.Channel, sentence)
			log.Printf("[*] %s", sentence)
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			continue
		}
	}

}
