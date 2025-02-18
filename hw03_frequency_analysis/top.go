package hw03frequencyanalysis

import (
	"slices"
	"strings"
)

func Top10(s string) []string {
	words := strings.Fields(s)
	wordsFreq := make(map[string]int)
	uniqWords := make([]string, 0, len(wordsFreq))

	for _, word := range words {
		wordsFreq[word]++
	}

	for word := range wordsFreq {
		uniqWords = append(uniqWords, word)
	}

	slices.SortFunc(uniqWords, func(a, b string) int {
		if wordsFreq[a] != wordsFreq[b] {
			return wordsFreq[b] - wordsFreq[a]
		}
		if a < b {
			return -1
		}
		return 1
	})

	if len(uniqWords) > 10 {
		uniqWords = uniqWords[:10]
	}
	return uniqWords
}
