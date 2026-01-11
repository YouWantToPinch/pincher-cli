package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
)

func handlerPayee(s *State, c *handlerContext) error {
	if val, ok := c.ctxValues["action"]; ok {
		switch val {
		case "add":
			return handlePayeeAdd(s, c)
		case "list":
			return handlePayeeList(s, c)
		case "update":
			return handlePayeeUpdate(s, c)
		case "delete":
			return handlePayeeDelete(s, c)
		default:
			return fmt.Errorf("action not implemented")
		}
	} else {
		return fmt.Errorf("action was not saved to context")
	}
}

func handlePayeeAdd(s *State, c *handlerContext) error {
	name, _ := c.args.pfx()

	c.args.trackOptArgs(&c.cmd, "notes")
	notes, _ := c.args.pfx()

	payeeCreated, err := s.Client.CreatePayee(name, notes)
	if err != nil {
		return err
	}
	if payeeCreated {
		fmt.Println("Payee " + name + " successfully created as user: " + s.Session.Username)
		fmt.Println("See it with: `payee list`")
		return nil
	} else {
		return fmt.Errorf("payee could not be created")
	}
}

func handlePayeeList(s *State, c *handlerContext) error {
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

	payees, err := s.Client.GetPayees(includeQuery)
	if err != nil {
		return err
	}
	if len(payees) == 0 {
		fmt.Printf("No payees found belonging to budget %s. \n", s.Client.ViewedBudget.Name)
		return nil
	}
	fmt.Printf("Payees under budget %s: \n", s.Client.ViewedBudget.Name)
	sort.Slice(payees, func(i, j int) bool {
		return payees[i].Name < payees[j].Name
	})
	const uuidLength = 36
	maxLenName := MaxOfStrings(ExtractStrings(payees, func(b client.Payee) string { return b.Name }))
	maxLenNotes := MaxOfStrings(ExtractStrings(payees, func(b client.Payee) string { return b.Notes }))
	fmt.Printf("  %-*s | %-*s | %s\n", maxLenName, "NAME", uuidLength, "ID", "NOTES")
	fmt.Printf("  %s-+-%s-+-%s\n", nDashes(maxLenName), nDashes(uuidLength), nDashes(maxLenNotes))
	for _, payee := range payees {
		fmt.Printf("  %-*s  %-*s   %s\n", maxLenName, payee.Name, uuidLength, payee.ID, payee.Notes)
	}

	return nil
}

func handlePayeeUpdate(s *State, c *handlerContext) error {
	payeeName, _ := c.args.pfx()

	payees, err := s.Client.GetPayees("")
	if err != nil {
		return err
	}
	payee, err := findPayeeByName(payeeName, payees)
	if err != nil {
		return err
	}

	c.args.trackOptArgs(&c.cmd, "name")
	payloadName, err := c.args.pfx()
	if err != nil {
		payloadName = payee.Name
	}
	c.args.trackOptArgs(&c.cmd, "notes")
	payloadNotes, err := c.args.pfx()
	if err != nil {
		payloadNotes = payee.Notes
	}

	err = s.Client.UpdatePayee(payee.ID.String(), payloadName, payloadNotes)
	if err != nil {
		return err
	}
	fmt.Println("Payee updated with new information")
	return nil
}

func handlePayeeDelete(s *State, c *handlerContext) error {
	name, _ := c.args.pfx()
	c.args.trackOptArgs(&c.cmd, "replacement")
	newPayeeName, _ := c.args.pfx()

	payees, err := s.Client.GetPayees("")
	if err != nil {
		return err
	}
	payee, err := findPayeeByName(name, payees)
	if err != nil {
		return err
	}

	err = s.Client.DeletePayee(payee.ID.String(), newPayeeName)
	if err != nil {
		return err
	}
	fmt.Println("Payee deleted.")
	return nil
}
