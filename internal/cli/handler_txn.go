package cli

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
)

func handlerTxn(s *State, c *handlerContext) error {
	if val, ok := c.ctxValues["action"]; ok {
		switch val {
		case "log":
			return handleTxnLog(s, c)
		case "transfer", "tfr":
			return handleTxnTransfer(s, c)
		case "list":
			return handleTxnList(s, c)
		default:
			return fmt.Errorf("ERROR: action not implemented")
		}
	} else {
		return fmt.Errorf("ERROR: action was not saved to context")
	}
}

func handleTxnTransfer(s *State, c *handlerContext) error {
	fromAccountName, _ := c.args.pfx()
	toAccountName, _ := c.args.pfx()
	amounts := map[string]int64{}
	{
		amount, _ := c.args.pfx()
		parsedAmount, err := parseCurrencyFromString(amount, s.Config.CurrencyISOCode)
		if err != nil {
			return fmt.Errorf("could not log transfer: %w", err)
		}
		// Pincher-CLI handles transfers in a deliberately from->to manner
		if parsedAmount < 0 {
			return fmt.Errorf("could not log transfer: amount to transfer must be positive")
		}
		amounts["TRANSFER"] = int64(parsedAmount) * -1
	}
	c.args.trackOptArgs(&c.cmd, "notes")
	notes, _ := c.args.pfx()
	c.args.trackOptArgs(&c.cmd, "date")
	transactionDate, err := c.args.pfx()
	if err != nil {
		transactionDate = time.Now().Format("2006-01-02")
	}
	c.args.trackOptArgs(&c.cmd, "cleared")
	isCleared, _ := c.args.pfx()

	_, err = s.Client.LogTxn(
		fromAccountName,
		toAccountName,
		transactionDate,
		"",
		notes,
		isCleared == "SET",
		amounts,
	)
	if err != nil {
		return err
	}
	return nil
}

func handleTxnLog(s *State, c *handlerContext) error {
	accountName, _ := c.args.pfx()
	payeeName, _ := c.args.pfx()
	totalAmountString, _ := c.args.pfx()
	category, _ := c.args.pfx()
	c.args.trackOptArgs(&c.cmd, "date")
	transactionDate, err := c.args.pfx()
	if err != nil {
		transactionDate = time.Now().Format("2006-01-02")
	}
	c.args.trackOptArgs(&c.cmd, "notes")
	notes, _ := c.args.pfx()
	c.args.trackOptArgs(&c.cmd, "cleared")
	isCleared, _ := c.args.pfx()

	amounts := map[string]int64{}
	c.args.trackOptArgs(&c.cmd, "split")
	splitArg, err := c.args.pfx()
	if err == nil {
		if strings.ToUpper(category) == "SPLIT" {

			splits := strings.Split(splitArg, ",")
			if len(splits) == 0 {
				return fmt.Errorf("split option used, but no splits provided")
			}
			var splitsTotal int64
			for _, split := range splits {
				pair := strings.Split(split, "=")
				if len(pair) != 2 {
					return fmt.Errorf("could not parse one or more splits")
				}
				category, amount := pair[0], pair[1]
				parsedAmount, err := parseCurrencyFromString(amount, s.Config.CurrencyISOCode)
				if err != nil {
					return err
				}
				amounts[category] = int64(parsedAmount)
				splitsTotal += int64(parsedAmount)
			}
			totalAmount, err := parseCurrencyFromString(totalAmountString, s.Config.CurrencyISOCode)
			if err != nil {
				return err
			}
			if splitsTotal != int64(totalAmount) {
				return fmt.Errorf("split amounts (%d) do not amount to total: %d", splitsTotal, totalAmount)
			}
		} else {
			return fmt.Errorf("substitute 'split' for the category argument to use the --splits option")
		}
	} else {
		totalAmount, err := parseCurrencyFromString(totalAmountString, s.Config.CurrencyISOCode)
		if err != nil {
			return err
		}
		amounts[category] = int64(totalAmount)
	}

	_, err = s.Client.LogTxn(
		accountName,
		"",
		transactionDate,
		payeeName,
		notes,
		isCleared == "SET",
		amounts,
	)
	if err != nil {
		return err
	}
	return nil
}

func handleTxnList(s *State, c *handlerContext) error {
	c.args.trackOptArgs(&c.cmd, "account")
	accountName, _ := c.args.pfx()
	c.args.trackOptArgs(&c.cmd, "category")
	categoryName, _ := c.args.pfx()
	c.args.trackOptArgs(&c.cmd, "payee")
	payeeName, _ := c.args.pfx()

	txnQuery := ""
	if len(accountName+categoryName+payeeName) > 0 {
		params := url.Values{}
		for k, v := range map[string]string{"account_name": accountName, "category_name": categoryName, "payee_name": payeeName} {
			if v != "" {
				params.Add(k, v)
			}
		}
		txnQuery = params.Encode()
		txnQuery = "?" + txnQuery
	}

	txns, err := s.Client.GetTxns(txnQuery)
	if err != nil {
		return err
	}
	if len(txns) == 0 {
		fmt.Printf("No transactions found under budget %s. \n", s.Client.ViewedBudget.Name)
		return nil
	}
	fmt.Printf("%s transactions: \n", s.Client.ViewedBudget.Name)
	sort.Slice(txns, func(i, j int) bool {
		return txns[i].TransactionDate.Before(txns[j].TransactionDate)
	})
	// const uuidLength = 36
	maxLenDate := MaxOfStrings(ExtractStrings(txns, func(t client.TransactionDetail) string { return t.TransactionDate.Format("2006-01-02") }))
	maxLenAmount := MaxOfStrings(ExtractStrings(txns, func(t client.TransactionDetail) string {
		return CurrencyUnit(t.TotalAmount).Format(s.Config.CurrencyISOCode, true)
	}))
	maxLenNotes := MaxOfStrings(ExtractStrings(txns, func(t client.TransactionDetail) string { return firstNChars(t.Notes, 25) + "..." }))
	fmt.Printf("  %-*s | %-*s | %s\n", maxLenDate, "DATE", maxLenAmount, "AMOUNT", "NOTES")
	fmt.Printf("  %s-+-%s-+-%s\n", nDashes(maxLenDate), nDashes(maxLenAmount), nDashes(maxLenNotes))
	for _, txn := range txns {
		fmt.Printf("  %-*s  %-*s   %s\n", maxLenDate, txn.TransactionDate.Format("2006-01-02"), maxLenAmount, CurrencyUnit(txn.TotalAmount).Format(s.Config.CurrencyISOCode, true), firstNChars(txn.Notes, 25))
	}

	return nil
}
