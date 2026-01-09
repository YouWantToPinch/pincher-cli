package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
)

func handlerAccount(s *State, c *handlerContext) error {
	if val, ok := c.ctxValues["action"]; ok {
		switch val {
		case "add":
			return handleAccountAdd(s, c)
		case "list":
			return handleAccountList(s, c)
		case "update":
			return handleAccountUpdate(s, c)
		case "restore":
			return handleAccountRestore(s, c)
		case "delete":
			return handleAccountDelete(s, c)
		default:
			return fmt.Errorf("action not implemented")
		}
	} else {
		return fmt.Errorf("action was not saved to context")
	}
}

func handleAccountAdd(s *State, c *handlerContext) error {
	name, _ := c.args.pfx()
	accountType, _ := c.args.pfx()

	c.args.trackOptArgs(&c.cmd, "notes")
	notes, _ := c.args.pfx()

	accountCreated, err := s.Client.CreateAccount(name, notes, accountType)
	if err != nil {
		return err
	}
	if accountCreated {
		fmt.Println("Account " + name + " successfully created as user: " + s.Session.Username + ".")
		fmt.Println("See it with: `account list`")
		return nil
	} else {
		return fmt.Errorf("account could not be created")
	}
}

func handleAccountList(s *State, c *handlerContext) error {
	c.args.trackOptArgs(&c.cmd, "include")
	includeQuery, _ := c.args.pfx()

	qualities := cleanInput(includeQuery)
	if len(qualities) > 0 {
		includeQuery = "?include="
		for i, quality := range qualities {
			includeQuery += strings.ToLower(quality)
			if i < len(qualities)-1 {
				includeQuery += "&"
			}
		}
	}

	accounts, err := s.Client.GetAccounts(includeQuery)
	if err != nil {
		return err
	}
	if len(accounts) == 0 {
		fmt.Printf("No accounts found belonging to budget %s. \n", s.Client.ViewedBudget.Name)
		return nil
	}
	fmt.Printf("Accounts under budget %s: \n", s.Client.ViewedBudget.Name)
	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].Name < accounts[j].Name
	})
	const uuidLength = 36
	maxLenName := MaxOfStrings(ExtractStrings(accounts, func(b client.Account) string { return b.Name }))
	maxLenNotes := MaxOfStrings(ExtractStrings(accounts, func(b client.Account) string { return b.Notes }))
	fmt.Printf("  %-*s | %-*s | %s\n", maxLenName, "NAME", uuidLength, "ID", "NOTES")
	fmt.Printf("  %s-+-%s-+-%s\n", nDashes(maxLenName), nDashes(uuidLength), nDashes(maxLenNotes))
	for _, account := range accounts {
		fmt.Printf("  %-*s  %-*s   %s\n", maxLenName, account.Name, uuidLength, account.ID, account.Notes)
	}

	return nil
}

func handleAccountUpdate(s *State, c *handlerContext) error {
	accountName, _ := c.args.pfx()

	accounts, err := s.Client.GetAccounts("")
	if err != nil {
		return err
	}
	account, err := findAccountByName(accountName, accounts)
	if err != nil {
		return err
	}

	c.args.trackOptArgs(&c.cmd, "name")
	payloadName, err := c.args.pfx()
	if err != nil {
		payloadName = account.Name
	}
	c.args.trackOptArgs(&c.cmd, "type")
	payloadAccountType, err := c.args.pfx()
	if err != nil {
		payloadAccountType = account.AccountType
	}
	c.args.trackOptArgs(&c.cmd, "notes")
	payloadNotes, err := c.args.pfx()
	if err != nil {
		payloadNotes = account.Notes
	}

	err = s.Client.UpdateAccount(account.ID.String(), payloadName, payloadNotes, payloadAccountType)
	if err != nil {
		return err
	}
	fmt.Println("Account updated with new information")
	return nil
}

func handleAccountRestore(s *State, c *handlerContext) error {
	accountName, _ := c.args.pfx()

	accounts, err := s.Client.GetAccounts("?include=deleted")
	if err != nil {
		return err
	}
	account, err := findAccountByName(accountName, accounts)
	if err != nil {
		return err
	}

	err = s.Client.RestoreAccount(account.ID.String(), accountName)
	if err != nil {
		return err
	}
	fmt.Println("Account restored")
	return nil
}

func handleAccountDelete(s *State, c *handlerContext) error {
	name, _ := c.args.pfx()
	c.args.trackOptArgs(&c.cmd, "hard")
	flagDeleteHard, _ := c.args.pfx()

	accounts, err := s.Client.GetAccounts("?include=deleted")
	if err != nil {
		return err
	}
	account, err := findAccountByName(name, accounts)
	if err != nil {
		return err
	}

	err = s.Client.DeleteAccount(account.ID.String(), flagDeleteHard)
	if err != nil {
		return err
	}
	return nil
}
