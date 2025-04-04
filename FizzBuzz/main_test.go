package main

import (
	"slices"
	"testing"
)

func TestFizzBuzz(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected []string
	}{
		{
			name:     "zero input",
			input:    0,
			expected: []string{},
		},
		{
			name:     "single number",
			input:    1,
			expected: []string{"1"},
		},
		{
			name:     "first Fizz",
			input:    3,
			expected: []string{"1", "2", "Fizz"},
		},
		{
			name:     "first Buzz",
			input:    5,
			expected: []string{"1", "2", "Fizz", "4", "Buzz"},
		},
		{
			name:     "first FizzBuzz",
			input:    15,
			expected: []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"},
		},
		{
			name:     "multiple FizzBuzz",
			input:    30,
			expected: []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz", "16", "17", "Fizz", "19", "Buzz", "Fizz", "22", "23", "Fizz", "Buzz", "26", "Fizz", "28", "29", "FizzBuzz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FizzBuzz(tt.input)
			if !slices.Equal(got, tt.expected) {
				t.Errorf("FizzBuzz(%d) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFizzBuzzNegative(t *testing.T) {
	got := FizzBuzz(-1)
	if len(got) != 0 {
		t.Errorf("FizzBuzz(-1) should return empty slice, got %v", got)
	}
}
