package main

import (
	"reflect"
	"testing"
)

func TestWordStats(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  WordSummary
	}{
		{
			name:  "basic case",
			input: "go go gophers love go",
			want: WordSummary{
				Counts:     map[string]int{"go": 3, "gophers": 1, "love": 1},
				MostCommon: "go",
				Frequency:  3,
			},
		},
		{
			name:  "empty string",
			input: "",
			want: WordSummary{
				Counts:     map[string]int{},
				MostCommon: "",
				Frequency:  0,
			},
		},
		{
			name:  "case insensitive",
			input: "Go gO Gophers love go",
			want: WordSummary{
				Counts:     map[string]int{"go": 3, "gophers": 1, "love": 1},
				MostCommon: "go",
				Frequency:  3,
			},
		},
		{
			name:  "punctuation",
			input: "Hello, world! Hello, world!",
			want: WordSummary{
				Counts:     map[string]int{"hello": 2, "world": 2},
				MostCommon: "hello",
				Frequency:  2,
			},
		},
		{
			name:  "multiple spaces",
			input: "  Hello   world ",
			want: WordSummary{
				Counts:     map[string]int{"hello": 1, "world": 1},
				MostCommon: "hello",
				Frequency:  1,
			},
		},
		{
			name:  "non-ASCII characters",
			input: "café café résumé résumé",
			want: WordSummary{
				Counts:     map[string]int{"café": 2, "résumé": 2},
				MostCommon: "café",
				Frequency:  2,
			},
		},
		{
			name:  "mixed punctuation",
			input: "Hello! How are you? I'm fine, thanks.",
			want: WordSummary{
				Counts:     map[string]int{"hello": 1, "how": 1, "are": 1, "you": 1, "im": 1, "fine": 1, "thanks": 1},
				MostCommon: "hello",
				Frequency:  1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WordStats(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WordStats(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// BenchmarkWordStats measures the performance of WordStats function
func BenchmarkWordStats(b *testing.B) {
	input := "This is a test string with multiple words and some punctuation! " +
		"Let's see how fast it can process this text. " +
		"Repeating words: test test test test test. " +
		"More punctuation: hello, world! how are you? I'm fine, thanks."

	for i := 0; i < b.N; i++ {
		WordStats(input)
	}
}
