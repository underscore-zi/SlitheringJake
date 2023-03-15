package markovchain

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type InsertionOpts struct {
	TerminateOnNewLine bool
}

// InsertLine assumes it is inserting a single terminated sequence. If there is no terminator it adds one,
// but if there is a terminator somewhere in the middle, it just treats it as part of the string
func (m *MarkovChain) InsertLine(line string) {
	line = strings.TrimSpace(line)
	if !m.isTerminated(line) {
		line += "."
	}

	terms := m.splitLine(line)
	if len(terms) <= m.n {
		// if the line is too short, just ignore it
		return
	}

	for i, term := range m.splitLine(line) {
		remaining := len(terms) - i
		if remaining <= m.n {
			continue
		}

		var builder strings.Builder
		builder.WriteString(term)

		endIndex := i + m.n
		if endIndex > len(terms)-1 {
			// if the endIndex reaches the end, then there is no continuation to add so just skip it
			continue
		}

		for _, xs := range terms[i+1 : endIndex] {
			builder.WriteRune(' ')
			builder.WriteString(xs)
		}

		start := builder.String()
		if m.isTerminated(start) {
			// terminator strings shouldn't be added to the front of the dictionary
			continue
		}

		builder.Reset()

		startIndex := endIndex // start of the continuation now
		endIndex += m.n

		// The > here and >= earlier is intentional, the first check if the endIndex is exactly the end I want to skip
		// since there is no continuation string, so the = isn't necessary here since we are consuming that last one
		if endIndex > len(terms)-1 {
			endIndex = len(terms)
		}

		for j, xs := range terms[startIndex:endIndex] {
			if j != 0 {
				builder.WriteRune(' ')
			}
			builder.WriteString(xs)
		}

		m.insertContinuation(start, builder.String())
		if i == 0 {
			// if this is the first term, then it is also a valid start
			m.insertContinuation("", start)
		}
	}
}

// insertContinuation is a utility that checks if the continuation already exists and if not, adds it
func (m *MarkovChain) insertContinuation(start, continuation string) {
	for _, existing := range m.dictionary[start] {
		if existing == continuation {
			return
		}
	}
	m.dictionary[start] = append(m.dictionary[start], continuation)
}
func (m *MarkovChain) InsertFile(filename string, opts InsertionOpts) error {
	// open a reader to the file
	reader, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", filename, err)
	}
	defer func() { _ = reader.Close() }()

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)

	builder := strings.Builder{}
	for scanner.Scan() {
		char := scanner.Text()
		r, _ := utf8.DecodeRuneInString(char)

		if unicode.IsSpace(r) {
			isNewLineTerminated := opts.TerminateOnNewLine && r == '\n'
			isNormalTerminated := m.isTerminated(builder.String())
			if isNewLineTerminated || isNormalTerminated {
				m.InsertLine(builder.String())
				builder.Reset()
				continue
			}
		}

		builder.WriteString(char)
	}
	if builder.String() != "" {
		m.InsertLine(builder.String())
	}
	return nil
}
