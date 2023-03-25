package main

import (
	"SlitheringJake/pkg/chatbot"
	"SlitheringJake/pkg/slitheringjake"
	"encoding/json"
	"flag"
	"log"
	"os"
)

var bot *chatbot.ChatBot

func main() {
	var err error
	var config slitheringjake.Config
	var configfn string

	flag.StringVar(&configfn, "config", "", "Path to config file (json)")
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

	jake, err := slitheringjake.NewSlitheringJake(config)
	if err != nil {
		log.Printf("Failed to create bot: %s", err.Error())
		return
	}

	err = jake.Bot.Run()
	if err != nil {
		panic(err)
	}

}
