package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

// Top returns the top 10 words. When taskWithAsteriskIsCompleted is true the
// extended rules are used: words are compared case-insensitively and
// punctuation at the beginning/end of a word is ignored.
func Top(x string, taskWithAsteriskIsCompleted bool) []string {
	if taskWithAsteriskIsCompleted {
		return Top10Asterisk(x)
	}

	return Top10(x)
}

// Top10 returns the 10 most frequent words, treating every whitespace
// separated token as a word.
func Top10(x string) []string {
	return Top10StringArr(strings.Fields(x))
}

// Top10Asterisk returns the 10 most frequent words using the extended rules.
func Top10Asterisk(s string) []string {
	return Top10StringArr(asteriskWords(s))
}

func asteriskWords(s string) []string {
	words := make([]string, 0)

	for _, token := range strings.Fields(strings.ToLower(s)) {
		word := strings.TrimFunc(token, unicode.IsPunct)
		if word == "" {
			// A token consisting only of punctuation is still a word
			// ("-------"), except for a single dash.
			if token == "-" {
				continue
			}
			word = token
		}
		words = append(words, word)
	}

	return words
}

// Top10StringArr returns the 10 most frequent values of words. Values with
// equal frequency are sorted lexicographically.
func Top10StringArr(words []string) []string {
	if len(words) == 0 {
		return []string{}
	}

	sorted := make([]string, len(words))
	copy(sorted, words)
	sort.Strings(sorted)

	type pair struct {
		key   string
		count int
	}

	frequencies := make([]pair, 0, len(sorted))
	start := 0

	for i := 1; i <= len(sorted); i++ {
		if i == len(sorted) || sorted[i] != sorted[start] {
			frequencies = append(frequencies, pair{key: sorted[start], count: i - start})
			start = i
		}
	}

	sort.Slice(frequencies, func(i, j int) bool {
		if frequencies[i].count == frequencies[j].count {
			return frequencies[i].key < frequencies[j].key
		}
		return frequencies[i].count > frequencies[j].count
	})

	limit := len(frequencies)
	if limit > 10 {
		limit = 10
	}

	result := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		result = append(result, frequencies[i].key)
	}

	return result
}
