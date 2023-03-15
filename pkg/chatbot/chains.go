package chatbot

import (
	"SlitheringJake/pkg/markovchain"
)

func (bot *ChatBot) NewChain(name, filename string, n int) {
	bot.lock("map_chains")
	defer bot.unlock("map_chains")

	newChain := markovchain.New(n)
	if err := newChain.InsertFile(filename, markovchain.InsertionOpts{TerminateOnNewLine: true}); err != nil {
		panic(err)
	}
	bot.newMutex("chain_" + name)
	bot.chains[name] = newChain
}

func (bot *ChatBot) GetChain(name string) *markovchain.MarkovChain {
	bot.lock("chain_" + name)
	return bot.chains[name]
}
func (bot *ChatBot) PutChain(name string) {
	bot.unlock("chain_" + name)
}
