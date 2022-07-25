package chatbot

import "SlitheringJake/pkg/mctg"

func (bot *ChatBot) NewChain(name, filename string, n int) {
	bot.lock("map_chains")
	defer bot.unlock("map_chains")

	newChain := mctg.New(n)
	newChain.LoadCorpus(filename, true)
	bot.newMutex("chain_" + name)
	bot.chains[name] = newChain
}

func (bot *ChatBot) GetChain(name string) *mctg.MCTG {
	bot.lock("chain_" + name)
	return bot.chains[name]
}
func (bot *ChatBot) PutChain(name string) {
	bot.unlock("chain_" + name)
}
