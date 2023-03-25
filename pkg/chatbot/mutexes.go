package chatbot

import "sync"

func (bot *ChatBot) NewMutex(name string) {
	bot.Lock("map_mutexes")
	defer bot.Unlock("map_mutexes")

	bot.mutexes[name] = &sync.Mutex{}

	// TODO: Check for duplicates and actually throw an error here
}

func (bot *ChatBot) Lock(name string) {
	if m, found := bot.mutexes[name]; found {
		m.Lock()
	} else {
		panic("Unknown mutex: " + name)
	}
}

func (bot *ChatBot) Unlock(name string) {
	if m, found := bot.mutexes[name]; found {
		m.Unlock()
	} else {
		panic("Unknown mutex: " + name)
	}
}
