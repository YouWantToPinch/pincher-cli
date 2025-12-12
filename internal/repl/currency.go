package repl

import (
	"strconv"
)

type Currency struct {
	ISOCode           string `json:"iso_code"`
	Symbol            string `json:"symbol"`
	DecimalSeparator  rune   `json:"decimal_separator"`
	ThousandSeparator rune   `json:"thousand_separator"`
	SymbolBefore      bool   `json:"symbol_before"`
	Name              string `json:"name"`
}

var Currencies = map[string]Currency{
	"USD": {"USD", "$", '.', ',', true, "US Dollar"},
	"CAD": {"CAD", "$", '.', ',', true, "Canadian Dollar"},
	"EUR": {"EUR", "€", ',', '.', true, "Euro"},
	"GBP": {"GBP", "£", '.', ',', true, "Pound Sterling"},
}

// CurrencyUnit is an abstract unit used to represent some amount in US cents.
type CurrencyUnit int64

// Format returns a string providing readers with the appropriate visual
// corresponding to their currency's ISO Code.
// NOTE: Further localization of this software may require
// a refactor of this formatting logic.
func (c CurrencyUnit) Format(ISOCode string, useSymbol bool) string {
	s := strconv.FormatInt(int64(c), 10)
	i := len(s) - 2
	res := s[:i] + string(Currencies[ISOCode].DecimalSeparator) + s[i:]
	if useSymbol {
		if Currencies[ISOCode].SymbolBefore {
			return Currencies[ISOCode].Symbol + res
		}
		return res + Currencies[ISOCode].Symbol
	}
	return res
}
