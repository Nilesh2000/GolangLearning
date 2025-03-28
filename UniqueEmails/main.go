package main

import (
	"strings"
)

func UniqueEmails(emails []string) int {
	seen := make(map[string]struct{})

	for _, email := range emails {
		parts := strings.Split(email, "@")
		local, domain := strings.ToLower(parts[0]), strings.ToLower(parts[1])

		local = strings.Split(local, "+")[0]
		local = strings.ReplaceAll(local, ".", "")

		normalized := local + "@" + domain
		seen[normalized] = struct{}{}
	}
	return len(seen)
}
