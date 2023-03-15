package markovchain

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixMilli())
}

// NormalizeToken a given word to a standard form
func (m *MarkovChain) NormalizeToken(word string) string {
	word = strings.ToLower(word)
	word = alphaNumeric.ReplaceAllString(word, "")
	return word
}

// splitLine turns a single line into a normalized split of terms
func (m *MarkovChain) splitLine(line string) []string {
	var out []string
	for _, word := range strings.Split(line, " ") {
		word = m.NormalizeToken(word)
		out = append(out, word)
	}
	return out
}

// isTerminated checks if the term ends with a terminator
func (m *MarkovChain) isTerminated(term string) bool {
	if len(term) == 0 {
		return false
	}
	lastCharacter := term[len(term)-1]
	for _, terminator := range terminators {
		if lastCharacter == uint8(terminator) {
			return true
		}
	}
	return false
}

// random returns a random number between 0 and max (exclusive)
func (m *MarkovChain) random(max int) int {
	return rand.Intn(max)
}

// randomStart returns a random starting term
func (m *MarkovChain) randomStart() string {
	starts := m.dictionary[""]
	if len(starts) == 0 {
		return ""
	}
	return starts[m.random(len(starts))]
}

// randomContinuation returns a random continuation for a given word
func (m *MarkovChain) randomContinuation(start string) string {
	continuations := m.dictionary[start]
	if len(continuations) == 0 {
		return ""
	}
	return continuations[m.random(len(continuations))]
}

func (m *MarkovChain) Print() {
	for k, v := range m.dictionary {
		fmt.Println("Key:", k)
		for _, xs := range v {
			fmt.Println("\t", xs)
		}
	}
}
