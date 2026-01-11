package cli

import (
	"fmt"
	"testing"
)

func Test_Format(t *testing.T) {
	tests := []struct {
		name      string
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

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s: %d", tt.isoCode, tt.input), func(t *testing.T) {
			displayStr := tt.input.Format(tt.isoCode, tt.useSymbol)
			if displayStr != tt.expected {
				t.Fatalf("expected string %s, but got string: %s", tt.expected, displayStr)
			}
		})
	}
}

func Test_ParseCurrencyFromString(t *testing.T) {
	tests := []struct {
		input    string
		isoCode  string
		expected CurrencyUnit
		wantErr  bool
	}{
		{
			input:    "200.00",
			isoCode:  "USD",
			expected: 20000,
		},
		{
			input:    "200,00",
			isoCode:  "EUR",
			expected: 20000,
		},
		{
			input:    "2,000",
			isoCode:  "USD",
			expected: 200000,
		},
		{
			input:    "2,000.00",
			isoCode:  "USD",
			expected: 200000,
		},
		{
			input:    "2,000,000.00",
			isoCode:  "USD",
			expected: 200000000,
		},
		{
			input:    "2.000.000,00",
			isoCode:  "EUR",
			expected: 200000000,
		},
		{
			input:    "501.39",
			isoCode:  "USD",
			expected: 50139,
		},
		{
			input:    "501,39",
			isoCode:  "EUR",
			expected: 50139,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s: %s", tt.isoCode, tt.input), func(t *testing.T) {
			amount, err := parseCurrencyFromString(tt.input, tt.isoCode)
			if tt.wantErr != (err != nil) {
				t.Fatalf("expected error: %v, but got: %v", tt.wantErr, (err != nil))
			}
			if amount != tt.expected {
				t.Fatalf("expected value %d, but got value: %d", tt.expected, amount)
			}
		})
	}
}
