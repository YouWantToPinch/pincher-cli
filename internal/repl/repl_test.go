package repl

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
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("ERROR: input vs expected are of unequal lengths")
			t.Fail()
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("ERROR: input word is unequal to expected word")
				fmt.Println("expected: ", expectedWord)
				fmt.Println("actual: ", word)

			}
		}
	}
}
