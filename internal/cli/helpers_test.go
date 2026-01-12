package cli

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " Hel lo  ",
			expected: []string{"Hel", "lo"},
		},
		{
			input:    "Hello, World!",
			expected: []string{"Hello,", "World!"},
		},
		{
			input:    "Hello World HELLO",
			expected: []string{"Hello", "World", "HELLO"},
		},
		{
			input:    "heLlO",
			expected: []string{"heLlO"},
		},
		{
			input:    `account add "My Checking Account" "on-budget" --notes "The checking account I use."`,
			expected: []string{"account", "add", "My Checking Account", "on-budget", "--notes", "The checking account I use."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			actual := cleanInput(tt.input)
			if len(actual) != len(tt.expected) {
				t.Errorf("input vs expected are of unequal lengths")
				t.Fail()
			}
			for i := range actual {
				phrase := actual[i]
				expectedPhrase := tt.expected[i]
				if phrase != expectedPhrase {
					t.Errorf("input does not match expected phrase")
				}
			}
		})
	}
}
