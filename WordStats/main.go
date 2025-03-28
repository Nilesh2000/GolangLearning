package main

import "strings"

type WordSummary struct {
	Counts     map[string]int
	MostCommon string
	Frequency  int
}

func WordStats(s string) WordSummary {
	dict := make(map[string]int)
	for _, word := range strings.Fields(s) {
		dict[strings.ToLower(word)] += 1
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
