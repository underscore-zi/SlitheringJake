package chatbot

import "sync"

func (bot *ChatBot) newMutex(name string) {
	bot.lock("map_mutexes")
	defer bot.unlock("map_mutexes")

	bot.mutexes[name] = &sync.Mutex{}

	// TODO: Check for duplicates and actually throw an error here
}

func (bot *ChatBot) lock(name string) {
	if m, found := bot.mutexes[name]; found {
		m.Lock()
	} else {
		panic("Unknown mutex: " + name)
	}
}

func (bot *ChatBot) unlock(name string) {
	if m, found := bot.mutexes[name]; found {
		m.Unlock()
	} else {
		panic("Unknown mutex: " + name)
	}
}
