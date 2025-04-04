package main

import (
	"fmt"
	"strings"
)

// UniqueEmails returns the number of unique email addresses after normalizing them.
// Normalization includes:
// - Converting to lowercase
// - Removing dots (.) from the local part
// - Removing everything after '+' in the local part
// Returns 0 if any email is invalid.
func UniqueEmails(emails []string) int {
	seen := make(map[string]struct{})

	for _, email := range emails {
		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			fmt.Printf("Invalid email format: %s\n", email)
			return 0
		}

		local, domain := strings.ToLower(parts[0]), strings.ToLower(parts[1])
		local = strings.Split(local, "+")[0]
		local = strings.ReplaceAll(local, ".", "")

		normalized := local + "@" + domain
		seen[normalized] = struct{}{}
	}
	return len(seen)
}
