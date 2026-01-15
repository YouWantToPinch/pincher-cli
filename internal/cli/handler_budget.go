package cli

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
	cc "github.com/YouWantToPinch/pincher-cli/internal/currency"
)

func handlerBudget(s *State, c *handlerContext) error {
	if val, ok := c.ctxValues["action"]; ok {
		switch val {
		case "add":
			return handleBudgetAdd(s, c)
		case "list":
			return handleBudgetList(s, c)
		case "report":
			return handleBudgetReport(s, c)
		case "view":
			return handleBudgetView(s, c)
		case "update":
			return handleBudgetUpdate(s, c)
		case "delete":
			return handleBudgetDelete(s, c)
		default:
			return fmt.Errorf("action not implemented")
		}
	} else {
		return fmt.Errorf("action was not saved to context")
	}
}

func handleBudgetAdd(s *State, c *handlerContext) error {
	name, _ := c.args.pfx()

	c.args.trackOptArgs(&c.cmd, "notes")
	notes, _ := c.args.pfx()

	budgetCreated, err := s.Client.CreateBudget(name, notes)
	if err != nil {
		return err
	}
	if budgetCreated {
		fmt.Println("Budget " + name + " successfully created as user: " + s.Session.Username + ".")
		fmt.Println("See it with: `budget view`")
		return nil
	} else {
		return fmt.Errorf("budget could not be created")
	}
}

func handleBudgetView(s *State, c *handlerContext) error {
	name, _ := c.args.pfx()

	budgets, err := s.Client.GetBudgets("")
	if err != nil {
		return fmt.Errorf("could not view specified budget: %w", err)
	}

	budget, err := findBudgetByName(name, budgets)
	if err != nil {
		return err
	}

	// NOTE: We store a VALUE rather than the ptr,
	// as the cache by nature may change at a moment's notice
	s.Client.ViewedBudget = *budget
	s.Session.OnViewBudget()
	fmt.Printf("Now viewing budget: %s\n", budget.Name)
	return nil
}

func handleBudgetReport(s *State, c *handlerContext) error {
	var err error

	monthTime := time.Now()
	c.args.trackOptArgs(&c.cmd, "month")
	month, _ := c.args.pfx()
	if month != "" {
		monthTime, err = time.Parse("2006-01", month)
		if err != nil {
			return fmt.Errorf("bad month format; use YYYY-MM")
		}
	}
	monthStr := monthTime.Format("2006-01-02")

	report, err := s.Client.GetBudgetReport(monthStr)
	if err != nil {
		return err
	}
	iso := s.Config.CurrencyISOCode
	assigned := cc.Format(report.Assigned, iso, true)
	activity := cc.Format(report.Activity, iso, true)
	balance := cc.Format(report.Balance, iso, true)
	fmt.Printf("%s report for %s:\n", monthStr, s.Client.ViewedBudget.Name)
	fmt.Printf("  %-*s | %-*s | %s\n", len(assigned), "ASSIGNED", len(activity), "ACTIVITY", "BALANCE")
	fmt.Printf("  %s-+-%s-+-%s\n", nDashes(len("ASSIGNED")), nDashes(len("ACTIVITY")), nDashes(len("BALANCE")))
	fmt.Printf("  %-*s | %-*s | %s\n", len(assigned), assigned, len(activity), activity, balance)
	fmt.Printf("  %s-+-%s-+-%s\n", nDashes(len(assigned)), nDashes(len(activity)), nDashes(len(balance)))
	return nil
}

func handleBudgetList(s *State, c *handlerContext) error {
	c.args.trackOptArgs(&c.cmd, "roles")
	roleQuery, _ := c.args.pfx()

	if roleQuery != "" {
		roles := cleanInput(roleQuery)
		roleQuery = "?roles="
		for i, role := range roles {
			roleQuery += strings.ToUpper(role)
			if i < len(roles)-1 {
				roleQuery += ","
			}
		}
	}

	budgets, err := s.Client.GetBudgets(roleQuery)
	if err != nil {
		return err
	}
	if len(budgets) == 0 {
		fmt.Printf("No memberships found in query from user %s. \n", s.Session.Username)
		return nil
	}

	fmt.Printf("%s's budget memberships: \n", s.Session.Username)
	sort.Slice(budgets, func(i, j int) bool {
		return budgets[i].Name < budgets[j].Name
	})
	const uuidLength = 36
	maxLenName := MaxOfStrings(ExtractStrings(budgets, func(b client.Budget) string { return b.Name })...)
	maxLenNotes := MaxOfStrings(ExtractStrings(budgets, func(b client.Budget) string { return b.Notes })...)
	fmt.Printf("  %-*s | %-*s | %s\n", maxLenName, "NAME", uuidLength, "ID", "NOTES")
	fmt.Printf("  %s-+-%s-+-%s\n", nDashes(maxLenName), nDashes(uuidLength), nDashes(maxLenNotes))
	for _, budget := range budgets {
		fmt.Printf("  %-*s | %-*s | %s\n", maxLenName, budget.Name, uuidLength, budget.ID, budget.Notes)
		fmt.Printf("  %s-+-%s-+-%s\n", nDashes(maxLenName), nDashes(uuidLength), nDashes(maxLenNotes))
	}

	return nil
}

func handleBudgetUpdate(s *State, c *handlerContext) error {
	budgetName, _ := c.args.pfx()

	budgets, err := s.Client.GetBudgets("")
	if err != nil {
		return err
	}
	budget, err := findBudgetByName(budgetName, budgets)
	if err != nil {
		return err
	}

	c.args.trackOptArgs(&c.cmd, "name")
	payloadName, err := c.args.pfx()
	if err != nil {
		payloadName = budget.Name
	}
	c.args.trackOptArgs(&c.cmd, "notes")
	payloadNotes, err := c.args.pfx()
	if err != nil {
		payloadNotes = budget.Notes
	}

	err = s.Client.UpdateBudget(budget.ID.String(), payloadName, payloadNotes)
	if err != nil {
		return err
	}
	fmt.Println("Budget info updated with new information")
	return nil
}

func handleBudgetDelete(s *State, c *handlerContext) error {
	name, _ := c.args.pfx()

	budgets, err := s.Client.GetBudgets("")
	if err != nil {
		return err
	}
	budget, err := findBudgetByName(name, budgets)
	if err != nil {
		return err
	}

	err = s.Client.DeleteBudget(budget.ID.String())
	if err != nil {
		return err
	}
	fmt.Println("Budget deleted.")
	return nil
}
