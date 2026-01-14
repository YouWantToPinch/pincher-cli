package cli

import "strings"

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
