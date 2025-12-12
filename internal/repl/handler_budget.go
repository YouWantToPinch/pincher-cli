package repl

import (
	"fmt"
	"sort"
	"strings"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
)

func handlerBudget(s *State, c handlerContext) error {
	action := c.args.pfx()
	switch action {
	case "add":
		return handleBudgetAdd(s, c)
	case "list":
		return handleBudgetList(s, c)
	case "view":
		// return handle_budgetView(s, c)
		return fmt.Errorf("ERROR: action not implemented")
	case "":
		return fmt.Errorf("ERROR: no action specified")
	default:
		return fmt.Errorf("ERROR: invalid action for budget: %s", action)
	}
}

func handleBudgetAdd(s *State, c handlerContext) error {
	name := c.args.pfx()
	c.args.trackOptArgs(&c.cmd, "notes")
	notes := c.args.pfx()

	budgetCreated, err := s.Client.CreateBudget(name, notes)
	if err != nil {
		return fmt.Errorf("ERROR: %s", err)
	}
	if budgetCreated {
		fmt.Println("Budget " + name + " successfully created as user: " + s.Client.LoggedInUser.Username + ".")
		fmt.Println("See it with: `budget view`")
		return nil
	} else {
		return fmt.Errorf("ERROR: budget could not be created")
	}
}

func handleBudgetList(s *State, c handlerContext) error {
	c.args.trackOptArgs(&c.cmd, "role")
	roleQuery := c.args.pfx()
	roles := cleanInput(roleQuery)
	if len(roles) > 0 {
		roleQuery = "?role="
		for i, role := range roles {
			roleQuery += strings.ToUpper(role)
			if i < len(roles)-1 {
				roleQuery += "&"
			}
		}
	}

	budgets, err := s.Client.GetBudgets(roleQuery)
	if err != nil {
		return fmt.Errorf("ERROR: %s", err)
	}
	if len(budgets) == 0 {
		fmt.Printf("No memberships found in query from user %s. \n", s.Client.LoggedInUser.Username)
		return nil
	}
	fmt.Printf("%s's budget memberships: \n", s.Client.LoggedInUser.Username)
	sort.Slice(budgets, func(i, j int) bool {
		return budgets[i].Name < budgets[j].Name
	})
	const uuidLength = 36
	maxLenName := MaxOfStrings(ExtractStrings(budgets, func(b client.Budget) string { return b.Name }))
	maxLenNotes := MaxOfStrings(ExtractStrings(budgets, func(b client.Budget) string { return b.Notes }))
	fmt.Printf("  %-*s | %-*s | %s\n", maxLenName, "NAME", uuidLength, "ID", "NOTES")
	fmt.Printf("  %s-+-%s-+-%s\n", nDashes(maxLenName), nDashes(uuidLength), nDashes(maxLenNotes))
	for _, budget := range budgets {
		fmt.Printf("  %-*s  %-*s   %s\n", maxLenName, budget.Name, uuidLength, budget.ID, budget.Notes)
	}

	return nil
}
