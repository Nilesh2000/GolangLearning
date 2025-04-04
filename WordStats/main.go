package main

import (
	"strings"
	"unicode"
)

// WordSummary represents the statistical analysis of words in a text.
// It contains a map of word counts, the most common word, and its frequency.
type WordSummary struct {
	Counts     map[string]int
	MostCommon string
	Frequency  int
}

// removePunctuation removes all punctuation characters from a word.
// It returns the cleaned word with all punctuation removed.
func removePunctuation(word string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return r
	}, word)
}

// WordStats analyzes a string and returns word statistics.
// It counts the occurrences of each word, finds the most common word,
// and calculates its frequency. Words are case-insensitive and punctuation is removed.
func WordStats(s string) WordSummary {
	dict := make(map[string]int)
	words := strings.Fields(s)

	// Pre-allocate map with estimated size to reduce reallocations
	if len(words) > 0 {
		dict = make(map[string]int, len(words))
	}

	for _, word := range words {
		cleanWord := removePunctuation(strings.ToLower(word))
		if cleanWord != "" {
			dict[cleanWord]++
		}
	}

	mostCommon := ""
	maxFrequency := 0
	for k, v := range dict {
		if v > maxFrequency {
			mostCommon = k
			maxFrequency = v
		}
	}

	return WordSummary{
		Counts:     dict,
		MostCommon: mostCommon,
		Frequency:  maxFrequency,
	}
}
