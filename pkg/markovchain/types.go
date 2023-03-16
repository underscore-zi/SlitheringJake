package markovchain

import (
	"regexp"
)

var alphaNumeric *regexp.Regexp
var terminators []rune

func init() {
	// terminators must be ASCII characters
	terminators = []rune{'.', '!', '?'}

	// alphaNumeric must include the terminators since they are needed to detect end of chains
	alphaNumeric = regexp.MustCompile("[^a-z0-9@-_.!?]")

	for _, terminator := range terminators {
		if terminator > 127 || terminator < 0 {
			panic("terminator must be a valid ASCII character")
		}
	}
}

type MarkovChain struct {
	dictionary map[string][]string
	sequences  []string
	n          int
}

func New(n int) *MarkovChain {
	return &MarkovChain{
		dictionary: make(map[string][]string),
		n:          n,
	}
}
