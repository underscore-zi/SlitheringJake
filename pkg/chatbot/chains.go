package chatbot

import (
	"SlitheringJake/pkg/markovchain"
)

// NewChain creates a new markov chain and adds it to the map of chains.
func (bot *ChatBot) NewChain(name, filename string, n int) {
	bot.Lock("map_chains")
	defer bot.Unlock("map_chains")

	newChain := markovchain.New(n)
	if err := newChain.InsertFile(filename, markovchain.InsertionOpts{TerminateOnNewLine: true}); err != nil {
		panic(err)
	}
	bot.NewMutex("chain_" + name)
	bot.chains[name] = newChain
}

// GetChain returns a pointer to the markov chain with the given name and locks access. This should always be paired with a call to PutChain.
func (bot *ChatBot) GetChain(name string) *markovchain.MarkovChain {
	bot.Lock("chain_" + name)
	return bot.chains[name]
}

// PutChain unlocks access to the markov chain with the given name. This should always be paired with a call to GetChain.
func (bot *ChatBot) PutChain(name string) {
	bot.Unlock("chain_" + name)
}
