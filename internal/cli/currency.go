package cli

import (
	"fmt"
	"strconv"
	"strings"
)

type Currency struct {
	ISOCode           string `json:"iso_code"`
	Symbol            string `json:"symbol"`
	DecimalSeparator  rune   `json:"decimal_separator"`
	ThousandSeparator rune   `json:"thousand_separator"`
	SymbolBefore      bool   `json:"symbol_before"`
	Name              string `json:"name"`
	DecimalFactor     int8   `json:"decimal_factor"`
}

var Currencies = map[string]Currency{
	"USD": {"USD", "$", '.', ',', true, "US Dollar", 100},
	"CAD": {"CAD", "$", '.', ',', true, "Canadian Dollar", 100},
	"EUR": {"EUR", "€", ',', '.', true, "Euro", 100},
	"GBP": {"GBP", "£", '.', ',', true, "Pound Sterling", 100},
}

// CurrencyUnit is an abstract unit used to represent
// some amount of money in the smallest unit of a
// currency.
type CurrencyUnit int64

// NOTE: Further localization of this software may
// warrant modification of this formatting and
// parsing logic.

// Format returns a string providing readers with the appropriate visual
// corresponding to their currency's ISO Code.
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

// parseCurrencyFromStr takes a string, attempts to parse
// it as the given currency by ISO, and if successful,
// converts the string to an integer representing the
// dollar unit of the currency, before multiplying the
// result by the factor necessary to produce the same
// value in the currency's smallest unit.
// If the string could could not be parsed as the given
// currency, an error is returned.
func parseCurrencyFromString(s string, currencyISO string) (CurrencyUnit, error) {
	currency, ok := Currencies[currencyISO]
	if !ok {
		return 0, fmt.Errorf("invalid or unavailable ISO Currency Code")
	}

	if s == "" {
		return 0, fmt.Errorf("no content to parse from input")
	}

	negateMult := 1
	if strings.HasPrefix(s, "-") {
		negateMult = -1
		s = strings.TrimLeft(s, "-")
	}

	if strings.HasPrefix(s, "0") {
		s = strings.TrimLeft(s, "0")
	}

	if strings.HasPrefix(s, string(currency.DecimalSeparator)) {
		s = "0" + s
	}

	var dollars int64
	var cents CurrencyUnit
	pair := strings.Split(s, string(currency.DecimalSeparator))
	switch len(pair) {
	case 2:
		if len(pair[1]) != 2 {
			return 0, fmt.Errorf("write any specified decimal values to only the hundredths place: .xy")
		}
		parsedCents, err := strconv.ParseInt(pair[1], 0, 64)
		if err != nil {
			return 0, fmt.Errorf("could not parse cent currency unit: %w", err)
		}
		cents = CurrencyUnit(parsedCents)
		fallthrough
	case 1:
		dollarString := pair[0]

		// WARN: Localization of a currency like Indian Rupees
		// would demand use of the Indian Numbering System, and
		// may call this modulo into question, necessitating the
		// aforementioned modification.

		// For the International System of Units,
		// separators for the portion preceding the
		// decimal should be found after every
		// third digit.
		if len(dollarString)%4 == 0 {
			return 0, fmt.Errorf("string could not be parsed as given currency")
		}
		// If the dollar amount is 1-999, dollarString may be parsed
		if len(dollarString) > 3 {
			dollarRunes := []rune(dollarString)
			i := 1
			for j := len(dollarRunes) - 1; j >= 0; j-- {
				if i%4 == 0 && dollarRunes[j] != currency.ThousandSeparator {
					return 0, fmt.Errorf("improper thousands separator for currency")
				}
				i++
			}
			dollarString = strings.ReplaceAll(dollarString, string(currency.ThousandSeparator), "")
		}
		parsedDollars, err := strconv.ParseInt(dollarString, 0, 64)
		if err != nil {
			return 0, fmt.Errorf("could not parse dollar unit: %w", err)
		}
		dollars = parsedDollars
	case 0:
		// sanity check; empty string check SHOULD keep from getting here
		return 0, fmt.Errorf("no content to parse from input")
	default:
		return 0, fmt.Errorf("decimal separator found more than once")
	}
	total := ((CurrencyUnit(dollars) * CurrencyUnit(currency.DecimalFactor)) + cents) * CurrencyUnit(negateMult)
	return total, nil
}
