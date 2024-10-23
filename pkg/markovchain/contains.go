package markovchain

import (
	"slices"
	"strings"
)

func (m *MarkovChain) generatePrefix(start string, iterations int) (string, float32) {
	if iterations > 10 {
		return "", 0.0
	}

	if _, found := m.dictionary[start]; !found {
		return "", 0.0
	}

	seen := make(map[string]bool)
	var attempts int
	var partsWithChoices int

	var parts []string

	for {
		prefix, success := m.randomPrefix(start)
		if !success {
			break
		}

		if xs, _ := m.prefixes(prefix); len(xs) > 1 {
			partsWithChoices++
		}

		if seen[prefix] {
			attempts++
			if attempts > 10 {
				break
			}
			continue
		}

		attempts = 0
		seen[prefix] = true
		if prefix == "" {
			break
		}
		start = prefix
		parts = append([]string{prefix}, parts...)
	}

	newLine := strings.Join(parts, " ")
	return newLine, float32(partsWithChoices) / float32(len(parts)+1)
}

func (m *MarkovChain) prefixes(token string) ([]string, bool) {
	var matches []string

	if _, found := m.dictionary[token]; !found {
		return nil, false
	}

	for key, entries := range m.dictionary {
		if slices.Contains(entries, token) {
			matches = append(matches, key)
		}
	}

	if len(matches) == 0 {
		return nil, false
	}

	return matches, true
}

func (m *MarkovChain) Contains(token string) (string, float32) {
	var matches []string
	for key, _ := range m.dictionary {
		if strings.Contains(" "+key, " "+token) {
			matches = append(matches, key)
		}
	}

	if len(matches) == 0 {
		return "", 0.0
	}

	var sentence string
	var score float32

	var attemptsRemaining = 10
	for attemptsRemaining >= 0 {
		attemptsRemaining--

		start := matches[m.random(len(matches))]
		prefix, prefixScore := m.generatePrefix(start, 0)
		cont, contScore := m.generateFrom(start, 0)

		if prefixScore == 0.0 || contScore == 0.0 {
			continue
		}

		sentence = prefix + " " + cont
		score = (prefixScore + contScore) / 2.0

		if m.isUnique(sentence) {
			break
		}
	}

	return sentence, score
}

func (m *MarkovChain) randomPrefix(token string) (string, bool) {
	prefixes, success := m.prefixes(token)
	if !success {
		return "", false
	}

	return prefixes[m.random(len(prefixes))], true
}
