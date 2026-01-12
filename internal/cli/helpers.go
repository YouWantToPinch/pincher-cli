package cli

import (
	"strings"
)

func ExtractStrings[T any](items []T, f func(T) string) []string {
	strings := make([]string, len(items))
	for i, v := range items {
		strings[i] = f(v)
	}
	return strings
}

func MaxOfStrings(s []string) int {
	maxLen := 0
	for _, str := range s {
		if len(str) > maxLen {
			maxLen = len(str)
		}
	}
	return maxLen
}

func cleanInput(text string) []string {
	fields := []string{}
	text = strings.TrimSpace(text)
	split := strings.Split(text, `"`)
	for i, substr := range split {
		if i%2 == 0 {
			addFields := strings.Fields(substr)
			fields = append(fields, addFields...)
		} else {
			fields = append(fields, substr)
		}
	}
	return fields
}

func nDashes(n int) string {
	return strings.Repeat("-", n)
}

// returns the first cmdElement with the given name from a slice of cmdElements
func findCMDElementWithName(elements []cmdElement, name string) (*cmdElement, bool) {
	for i := range elements {
		el := &elements[i]
		if el.name == name {
			return el, true
		}
	}
	return nil, false
}

func firstNChars(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
