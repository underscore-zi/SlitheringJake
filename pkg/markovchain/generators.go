package markovchain

import (
	"strings"
)

func (m *MarkovChain) generateFrom(start string, iterations int) (string, float32) {
	if iterations > 10 {
		return "", 0.0
	}
	if _, found := m.dictionary[start]; !found {
		return "", 0.0
	}
	originalStart := start

	var builder strings.Builder
	seen := make(map[string]bool)
	builder.WriteString(start)

	var attempts, parts, partsWithChoices int

	for {
		continuation := m.randomContinuation(start)

		parts++
		if len(m.dictionary[start]) > 1 {
			partsWithChoices++
		}

		if continuation == "" {
			break
		}

		// Try to avoid any loops and repetitions
		if seen[continuation] {
			attempts++
			if attempts > 10 {
				break
			}
			continue
		}
		attempts = 0
		seen[continuation] = true

		builder.WriteRune(' ')
		builder.WriteString(continuation)
		start = continuation
		if m.isTerminated(continuation) {
			break
		}
	}

	newLine := strings.TrimSpace(builder.String())
	if m.isUnique(newLine) {
		return newLine, float32(partsWithChoices) / float32(parts)
	} else {
		return m.generateFrom(originalStart, iterations+1)
	}

}

// Generate returns a random string and a float representing the average number of options per token
func (m *MarkovChain) Generate() (string, float32) {
	for {
		start := m.randomStart()
		if start == "" {
			return "", 0.0
		}
		gen, uniq := m.generateFrom(start, 0)
		if gen == "" {
			continue
		} else {
			return gen, uniq
		}
	}
}

// StartsWith returns a random string but the first token should contain the desired prefix
func (m *MarkovChain) StartsWith(prefix string) (string, float32) {
	// Despite the name StartsWith it is acceptable for the prefix to be a substring of the start
	var starters []string
	for _, starts := range m.dictionary[""] {
		if strings.Contains(starts, prefix) {
			starters = append(starters, starts)
		}
	}
	if len(starters) == 0 {
		return "", 0.0
	}

	for i := 0; i < 10; i++ {
		start := starters[m.random(len(starters))]
		gen, uniq := m.generateFrom(start, 0)
		if gen == "" {
			continue
		} else {
			return gen, uniq
		}
	}
	return "", 0.0
}
