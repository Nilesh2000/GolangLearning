package main

import "testing"

func TestUniqueEmails(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  int
	}{
		{
			name: "emails with dots and pluses",
			input: []string{
				"test.email+spam@gmail.com",
				"test.e.mail@gmail.com",
				"testemail@gmail.com",
				"user@domain.com",
			},
			want: 2,
		},
		{
			name: "duplicate emails with different casing or spacing",
			input: []string{
				"Alice.Z+promo@domain.com",
				"alicez@domain.com",
				"alice.z+foo@domain.com",
				"bob@domain.com",
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UniqueEmails(tt.input)
			if got != tt.want {
				t.Errorf("got %q want %q", got, tt.want)
			}
		})
	}
}
