package repl

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
	split := strings.Split(text, `"`)
	for i, substr := range split {
		if i%2 == 0 {
			lower := strings.ToLower(substr)
			addFields := strings.Fields(lower)
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
