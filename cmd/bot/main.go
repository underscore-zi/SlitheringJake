package main

import (
	"SlitheringJake/pkg/chatbot"
	"encoding/json"
	"flag"
	"log"
	"os"
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

	if content, err := os.ReadFile(configfn); err == nil {
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
			bot.NewChain("quotes", quotesfn, 2)
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
