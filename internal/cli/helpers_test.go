package cli

import (
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " Hel lo  ",
			expected: []string{"hel", "lo"},
		},
		{
			input:    "Hello, World!",
			expected: []string{"hello,", "world!"},
		},
		{
			input:    "Hello World HELLO",
			expected: []string{"hello", "world", "hello"},
		},
		{
			input:    "heLlO",
			expected: []string{"hello"},
		},
		{
			input:    `account add "My Checking Account" "on-budget" --notes "The checking account I use."`,
			expected: []string{"account", "add", "My Checking Account", "on-budget", "--notes", "The checking account I use."},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("input vs expected are of unequal lengths")
			t.Fail()
		}
		for i := range actual {
			phrase := actual[i]
			expectedPhrase := c.expected[i]
			if phrase != expectedPhrase {
				t.Errorf("input word is unequal to expected phrase")
				fmt.Println("expected: ", expectedPhrase)
				fmt.Println("actual: ", phrase)

			}
		}
	}
}
