package cli

import (
	"fmt"
	"sort"

	pgo "github.com/YouWantToPinch/pincher-sdk-go/pinchergo"
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
	c.args.trackOptArgs(&c.cmd, "off-budget")
	accountType, _ := c.args.pfx()
	if accountType == "SET" {
		accountType = pgo.BudgetAccountTypeOffBudget
	} else {
		accountType = pgo.BudgetAccountTypeOnBudget
	}

	c.args.trackOptArgs(&c.cmd, "notes")
	notes, _ := c.args.pfx()

	err := s.Client.BudgetAccountCreate(s.Session.ActiveBudget.ID.String(), pgo.BudgetAccountCreateData{
		MetaData: pgo.MetaData{
			Name:  name,
			Notes: notes,
		},
		AccountType: accountType,
	})
	if err != nil {
		return fmt.Errorf("s.Client.BudgetAccountCreate: %w", err)
	} else {
		fmt.Println("Account " + name + " successfully created as user: " + s.Session.ActiveUser.Username + ".")
		fmt.Println("See it with: `account list`")
		return nil
	}
}

func handleAccountList(s *State, c *handlerContext) error {
	c.args.trackOptArgs(&c.cmd, "deleted")
	listDeleted, _ := c.args.pfx()

	listDeletedQuery := ""
	if listDeleted == "SET" {
		listDeletedQuery = "?deleted"
	}

	accounts, err := s.GetAccounts(s.Session.ActiveBudget.ID.String(), listDeletedQuery)
	if err != nil {
		return err
	}
	if len(accounts) == 0 {
		d := ""
		if listDeletedQuery != "" {
			d = "deleted "
		}
		fmt.Printf("No %saccounts found belonging to budget %s. \n", d, s.Session.ActiveBudget.Name)
		return nil
	}
	fmt.Printf("Accounts under budget %s: \n", s.Session.ActiveBudget.Name)
	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].Name < accounts[j].Name
	})
	const uuidLength = 36
	maxLenName := MaxOfStrings(ExtractStrings(accounts, func(b *pgo.Account) string { return b.Name })...)
	maxLenNotes := MaxOfStrings(ExtractStrings(accounts, func(b *pgo.Account) string { return b.Notes })...)
	fmt.Printf("  %-*s | %-*s | %s\n", maxLenName, "NAME", uuidLength, "ID", "NOTES")
	fmt.Printf("  %s-+-%s-+-%s\n", nDashes(maxLenName), nDashes(uuidLength), nDashes(maxLenNotes))
	for _, account := range accounts {
		fmt.Printf("  %-*s  %-*s   %s\n", maxLenName, account.Name, uuidLength, account.ID, account.Notes)
	}

	return nil
}

func handleAccountUpdate(s *State, c *handlerContext) error {
	accountName, _ := c.args.pfx()

	accounts, err := s.GetAccounts(s.Session.ActiveBudget.ID.String(), "")
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

	err = s.Client.BudgetAccountUpdate(s.Session.ActiveBudget.ID.String(), account.ID.String(), pgo.BudgetAccountUpdateData{
		MetaData: pgo.MetaData{
			Name:  payloadName,
			Notes: payloadNotes,
		},
		AccountType: payloadAccountType,
	})
	if err != nil {
		return err
	}
	fmt.Println("Account updated with new information")
	return nil
}

func handleAccountRestore(s *State, c *handlerContext) error {
	name, _ := c.args.pfx()

	accounts, err := s.GetAccounts(s.Session.ActiveBudget.ID.String(), "?deleted")
	if err != nil {
		return err
	}
	account, err := findAccountByName(name, accounts)
	if err != nil {
		return err
	}

	err = s.Client.BudgetAccountRestore(s.Session.ActiveBudget.ID.String(), account.ID.String())
	if err != nil {
		return err
	}
	fmt.Println("Account restored.")
	return nil
}

func handleAccountDelete(s *State, c *handlerContext) error {
	name, _ := c.args.pfx()
	c.args.trackOptArgs(&c.cmd, "hard")
	flagDeleteHard, _ := c.args.pfx()

	deleteHard := flagDeleteHard == "SET"
	deleteQuery := ""
	if deleteHard {
		deleteQuery = "?deleted"
	}

	accounts, err := s.GetAccounts(s.Session.ActiveBudget.ID.String(), deleteQuery)
	if err != nil {
		return err
	}
	account, err := findAccountByName(name, accounts)
	if err != nil {
		return err
	}

	err = s.Client.BudgetAccountDelete(s.Session.ActiveBudget.ID.String(), account.ID.String(), pgo.BudgetAccountDeleteData{
		DeleteHard: deleteHard,
	})
	if err != nil {
		return err
	}
	if deleteHard {
		fmt.Println("Account deleted. It cannot be restored.")
	} else {
		fmt.Println("Account is deleted. It may be restored, or permanently deleted.")
	}
	return nil
}
