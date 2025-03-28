package main

import (
	"reflect"
	"testing"
)

func TestWordStats(t *testing.T) {
	tests := []struct {
		input string
		want  WordSummary
	}{
		{
			input: "go go gophers love go",
			want: WordSummary{
				Counts:     map[string]int{"go": 3, "gophers": 1, "love": 1},
				MostCommon: "go",
				Frequency:  3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := WordStats(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %q want %q", got, tt.want)
			}
		})
	}
}
