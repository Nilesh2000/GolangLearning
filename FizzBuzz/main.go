package main

import "strconv"

const (
	fizz     = 3
	buzz     = 5
	fizzBuzz = fizz * buzz
)

func FizzBuzz(n int) []string {
	if n <= 0 {
		return []string{}
	}

	output := make([]string, n)
	for i := 1; i <= n; i++ {
		switch {
		case i%fizzBuzz == 0:
			output[i-1] = "FizzBuzz"
		case i%fizz == 0:
			output[i-1] = "Fizz"
		case i%buzz == 0:
			output[i-1] = "Buzz"
		default:
			output[i-1] = strconv.Itoa(i)
		}
	}
	return output
}
