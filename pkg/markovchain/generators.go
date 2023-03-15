package markovchain

import "strings"

func (m *MarkovChain) generateFrom(start string) (string, float32) {
	if _, found := m.dictionary[start]; !found {
		return "", 0.0
	}

	var builder strings.Builder
	seen := make(map[string]bool)
	builder.WriteString(start)

	var attempts, options, parts int

	for {
		continuation := m.randomContinuation(start)

		options += len(m.dictionary[start])
		parts++

		if continuation == "" {
			break
			//continuation = m.Generate()
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

	trimmed := strings.TrimSpace(builder.String())
	return trimmed, float32(options) / float32(parts)

}

// Generate returns a random string and a float representing the average number of options per token
func (m *MarkovChain) Generate() (string, float32) {
	start := m.randomStart()
	if start == "" {
		return "", 0.0
	}
	return m.generateFrom(start)
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

	start := starters[m.random(len(starters))]
	return m.generateFrom(start)
}
