package cli

import (
	"os"
	"strings"

	"golang.org/x/term"
)

func ExtractStrings[T any](items []T, f func(T) string) []string {
	strings := make([]string, len(items))
	for i, v := range items {
		strings[i] = f(v)
	}
	return strings
}

func MaxOfStrings(s ...string) int {
	maxLen := 0
	for _, str := range s {
		if len(str) > maxLen {
			maxLen = len(str)
		}
	}
	return maxLen
}

func nDashes(n int) string {
	return strings.Repeat("-", n)
}

func firstNChars(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}
	return width
}

func wrapText(s string, width int, indent int) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}

	lineLen := 0
	var out strings.Builder

	for i, w := range words {
		if i == 0 {
			out.WriteString(w)
			lineLen = len(w)
			continue
		}

		if lineLen+1+len(w) > width {
			out.WriteString("\n")
			out.WriteString(strings.Repeat(" ", indent))
			out.WriteString(w)
			lineLen = len(w)
		} else {
			out.WriteString(" ")
			out.WriteString(w)
			lineLen += 1 + len(w)
		}
	}

	return out.String()
}

func makeAlignedTable(column1 []string, column2 []string) string {
	if len(column1) == 0 || len(column1) != len(column2) {
		return ""
	}

	termWidth := max(getTerminalWidth(), 20)

	// Get max length of column 1 for alignment purposes
	maxKeyLen := 0
	for _, k := range column1 {
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	var out strings.Builder
	out.Grow(len(column1) * termWidth)

	for i, key := range column1 {
		val := column2[i]

		out.WriteString(key)
		if pad := maxKeyLen - len(key); pad > 0 {
			out.WriteString(strings.Repeat(" ", pad))
		}
		out.WriteString("  ")

		prefixLen := maxKeyLen + 2
		availWidth := termWidth - prefixLen
		availWidth = max(availWidth, 1)

		wrapped := wrapText(val, availWidth, prefixLen)
		out.WriteString(wrapped)
		out.WriteByte('\n')
	}

	return out.String()
}
