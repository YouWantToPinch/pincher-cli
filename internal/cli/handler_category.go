package cli

import (
	"fmt"
	"sort"
	"time"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
	cc "github.com/YouWantToPinch/pincher-cli/internal/currency"
)

func handlerCategory(s *State, c *handlerContext) error {
	if val, ok := c.ctxValues["action"]; ok {
		switch val {
		case "add":
			return handleCategoryAdd(s, c)
		case "list":
			return handleCategoryList(s, c)
		case "update":
			return handleCategoryUpdate(s, c)
		case "delete":
			return handleCategoryDelete(s, c)
		case "assign":
			return handleCategoryAssign(s, c)
		case "reports":
			return handleCategoryReports(s, c)
		default:
			return fmt.Errorf("action not implemented")
		}
	} else {
		return fmt.Errorf("action was not saved to context")
	}
}

func handleCategoryAdd(s *State, c *handlerContext) error {
	name, _ := c.args.pfx()

	c.args.trackOptArgs(&c.cmd, "notes")
	notes, _ := c.args.pfx()

	c.args.trackOptArgs(&c.cmd, "group")
	groupName, _ := c.args.pfx()

	err := s.Client.BudgetCategoryCreate(s.Session.ActiveBudget.ID.String(), client.BudgetCategoryCreateData{
		MetaData: client.MetaData{
			Name:  name,
			Notes: notes,
		},
		GroupName: groupName,
	})
	if err != nil {
		return err
	}
	fmt.Println("Category " + name + " successfully created as user: " + s.Session.ActiveUser.Username)
	fmt.Println("See it with: `category list`")
	return nil
}

func handleCategoryAssign(s *State, c *handlerContext) error {
	toCategory, _ := c.args.pfx()

	amount, _ := c.args.pfx()
	parsedAmount, err := cc.Parse(amount, s.Config.CurrencyISOCode)
	if err != nil {
		return err
	}

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

	c.args.trackOptArgs(&c.cmd, "from")
	fromCategory, _ := c.args.pfx()

	err = s.Client.BudgetCategoryAssign(s.Session.ActiveBudget.ID.String(), monthStr, client.BudgetCategoryAssignData{
		Amount:       parsedAmount,
		ToCategory:   toCategory,
		FromCategory: fromCategory,
	})
	if err != nil {
		return err
	}
	if fromCategory == "" {
		fmt.Printf("Assigned %s to category %s for month %s\n", amount, toCategory, monthStr)
	} else {
		fmt.Printf("Assigned %s to category %s from %s in month %s\n", amount, toCategory, fromCategory, monthStr)
	}

	return nil
}

func handleCategoryReports(s *State, c *handlerContext) error {
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

	reports, err := s.Client.BudgetCategoryReports(s.Session.ActiveBudget.ID.String(), monthTime.Format("2006-01-02"))
	if err != nil {
		return err
	}
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Name < reports[j].Name
	})
	if len(reports) == 0 {
		fmt.Println("Nothing to report.")
		return nil
	}
	fmt.Printf("Categories under budget %s: \n", s.Session.ActiveBudget.Name)
	maxLenMonth := MaxOfStrings(ExtractStrings(reports, func(r *client.CategoryReport) string { return r.MonthID.Format("2006-01") })...)
	maxLenName := MaxOfStrings(ExtractStrings(reports, func(r *client.CategoryReport) string { return r.Name })...)
	maxLenAssigned := MaxOfStrings(ExtractStrings(reports, func(r *client.CategoryReport) string { return cc.Format(r.Assigned, s.Config.CurrencyISOCode, true) })...)
	maxLenActivity := MaxOfStrings(ExtractStrings(reports, func(r *client.CategoryReport) string { return cc.Format(r.Activity, s.Config.CurrencyISOCode, true) })...)
	maxLenBalance := MaxOfStrings(ExtractStrings(reports, func(r *client.CategoryReport) string { return cc.Format(r.Balance, s.Config.CurrencyISOCode, true) })...)
	fmt.Printf("  %-*s | %-*s | %-*s | %-*s | %s\n", maxLenMonth, "MONTH", maxLenName, "NAME", maxLenAssigned, "ASSIGNED", maxLenActivity, "ACTIVITY", "BALANCE")
	fmt.Printf("  %s-+-%s-+-%s-+-%s-+-%s\n", nDashes(maxLenMonth), nDashes(maxLenName), nDashes(maxLenAssigned), nDashes(maxLenActivity), nDashes(maxLenBalance))
	for _, report := range reports {
		fmt.Printf("  %-*s | %-*s | %-*s | %-*s | %s\n",
			maxLenMonth,
			report.MonthID.Format("2006-01"),
			maxLenName,
			report.Name,
			maxLenAssigned,
			cc.Format(report.Assigned, s.Config.CurrencyISOCode, true),
			maxLenActivity,
			cc.Format(report.Activity, s.Config.CurrencyISOCode, true),
			cc.Format(report.Balance, s.Config.CurrencyISOCode, true),
		)
	}

	return nil
}

func handleCategoryList(s *State, c *handlerContext) error {
	groupQuery := ""
	c.args.trackOptArgs(&c.cmd, "group")
	groupName, _ := c.args.pfx()
	if groupName != "" {
		groupQuery = "?group_name=" + groupName
	}

	categories, err := s.GetCategories(s.Session.ActiveBudget.ID.String(), groupQuery)
	if err != nil {
		return err
	}
	if len(categories) == 0 {
		fmt.Printf("No categories found belonging to budget %s. \n", s.Session.ActiveBudget.Name)
		return nil
	}
	fmt.Printf("Categories under budget %s: \n", s.Session.ActiveBudget.Name)
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Name < categories[j].Name
	})
	const uuidLength = 36
	maxLenName := MaxOfStrings(ExtractStrings(categories, func(b *client.Category) string { return b.Name })...)
	maxLenNotes := MaxOfStrings(ExtractStrings(categories, func(b *client.Category) string { return b.Notes })...)
	fmt.Printf("  %-*s | %-*s | %s\n", maxLenName, "NAME", uuidLength, "ID", "NOTES")
	fmt.Printf("  %s-+-%s-+-%s\n", nDashes(maxLenName), nDashes(uuidLength), nDashes(maxLenNotes))
	for _, category := range categories {
		fmt.Printf("  %-*s  %-*s  %s\n", maxLenName, category.Name, uuidLength, category.ID, category.Notes)
	}

	return nil
}

func handleCategoryUpdate(s *State, c *handlerContext) error {
	categoryName, _ := c.args.pfx()

	categories, err := s.GetCategories(s.Session.ActiveBudget.ID.String(), "")
	if err != nil {
		return err
	}
	category, err := findCategoryByName(categoryName, categories)
	if err != nil {
		return err
	}

	c.args.trackOptArgs(&c.cmd, "group")
	groupName, _ := c.args.pfx()

	c.args.trackOptArgs(&c.cmd, "name")
	payloadName, err := c.args.pfx()
	if err != nil {
		payloadName = category.Name
	}
	c.args.trackOptArgs(&c.cmd, "notes")
	payloadNotes, err := c.args.pfx()
	if err != nil {
		payloadNotes = category.Notes
	}

	err = s.Client.BudgetCategoryUpdate(s.Session.ActiveBudget.ID.String(), category.ID.String(), client.BudgetCategoryUpdateData{
		MetaData: client.MetaData{
			Name:  payloadName,
			Notes: payloadNotes,
		},
		GroupName: groupName,
	})
	if err != nil {
		return err
	}
	fmt.Println("Category updated with new information")
	return nil
}

func handleCategoryDelete(s *State, c *handlerContext) error {
	name, _ := c.args.pfx()

	categories, err := s.GetCategories(s.Session.ActiveBudget.ID.String(), "")
	if err != nil {
		return err
	}
	category, err := findCategoryByName(name, categories)
	if err != nil {
		return err
	}

	err = s.Client.BudgetCategoryDelete(s.Session.ActiveBudget.ID.String(), category.ID.String())
	if err != nil {
		return err
	}
	fmt.Println("Category deleted.")
	return nil
}
