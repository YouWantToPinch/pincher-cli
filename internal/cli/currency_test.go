package cli

import (
	"testing"
)

func Test_ToString(t *testing.T) {
	cases := []struct {
		input     CurrencyUnit
		isoCode   string
		useSymbol bool
		expected  string
	}{
		{
			input:     20000,
			isoCode:   "USD",
			useSymbol: false,
			expected:  "200.00",
		},
		{
			input:     20000,
			isoCode:   "USD",
			useSymbol: true,
			expected:  "$200.00",
		},
		{
			input:     50000 + 139,
			isoCode:   "USD",
			useSymbol: false,
			expected:  "501.39",
		},
		{
			input:     42000 + 69,
			isoCode:   "USD",
			useSymbol: false,
			expected:  "420.69",
		},
	}

	for _, c := range cases {
		displayStr := c.input.Format(c.isoCode, c.useSymbol)
		if displayStr != c.expected {
			t.Errorf("ERROR: expected string %s, but got string: %s", c.expected, displayStr)
			t.Fail()
		}
	}
}
