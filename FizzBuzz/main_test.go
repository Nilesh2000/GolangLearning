package main

import (
	"slices"
	"testing"
)

func TestFizzBuzz(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  []string
	}{
		{
			name:  "n is 0",
			input: 0,
			want:  []string{},
		},
		{
			name:  "n is 1",
			input: 1,
			want:  []string{"1"},
		},
		{
			name:  "n is 5",
			input: 5,
			want:  []string{"1", "2", "Fizz", "4", "Buzz"},
		},
		{
			name:  "n is 15",
			input: 15,
			want:  []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FizzBuzz(tt.input)
			if !slices.Equal(got, tt.want) {
				t.Errorf("got %q want %q", got, tt.want)
			}
		})
	}
}
