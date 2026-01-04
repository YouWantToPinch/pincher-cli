package repl

import (
	"fmt"
	"sort"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
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
	assignGroupID := ""
	if groupName != "" {
		groups, err := s.Client.GetGroups("")
		if err != nil {
			return err
		}
		group, err := findGroupByName(groupName, groups)
		if err != nil {
			return err
		}
		assignGroupID = group.ID.String()
	}

	categoryCreated, err := s.Client.CreateCategory(name, notes, assignGroupID)
	if err != nil {
		return err
	}
	if categoryCreated {
		fmt.Println("Category " + name + " successfully created as user: " + s.Client.LoggedInUser.Username + ".")
		fmt.Println("See it with: `category list`")
		return nil
	} else {
		return fmt.Errorf("budget could not be created")
	}
}

func handleCategoryList(s *State, c *handlerContext) error {
	groupQuery := ""
	c.args.trackOptArgs(&c.cmd, "group")
	groupName, _ := c.args.pfx()
	searchGroupID := ""
	if groupName != "" {
		groups, err := s.Client.GetGroups("")
		if err != nil {
			return err
		}
		group, err := findGroupByName(groupName, groups)
		if err != nil {
			return err
		}
		searchGroupID = group.ID.String()
		groupQuery += "?group_id="
	}
	groupQuery += searchGroupID

	categories, err := s.Client.GetCategories(groupQuery)
	if err != nil {
		return err
	}
	if len(categories) == 0 {
		fmt.Printf("No categories found belonging to budget %s. \n", s.Client.ViewedBudget.Name)
		return nil
	}
	fmt.Printf("Categories under budget %s: \n", s.Client.ViewedBudget.Name)
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Name < categories[j].Name
	})
	const uuidLength = 36
	maxLenName := MaxOfStrings(ExtractStrings(categories, func(b client.Category) string { return b.Name }))
	maxLenNotes := MaxOfStrings(ExtractStrings(categories, func(b client.Category) string { return b.Notes }))
	fmt.Printf("  %-*s | %-*s | %s\n", maxLenName, "NAME", uuidLength, "ID", "NOTES")
	fmt.Printf("  %s-+-%s-+-%s\n", nDashes(maxLenName), nDashes(uuidLength), nDashes(maxLenNotes))
	for _, category := range categories {
		fmt.Printf("  %-*s  %-*s   %s\n", maxLenName, category.Name, uuidLength, category.ID, category.Notes)
	}

	return nil
}

func handleCategoryUpdate(s *State, c *handlerContext) error {
	categoryName, _ := c.args.pfx()

	categories, err := s.Client.GetCategories("")
	if err != nil {
		return err
	}
	category, err := findCategoryByName(categoryName, categories)
	if err != nil {
		return err
	}

	c.args.trackOptArgs(&c.cmd, "group")
	groupName, _ := c.args.pfx()
	assignGroupID := ""
	if groupName != "" {
		groups, err := s.Client.GetGroups("")
		if err != nil {
			return err
		}
		group, err := findGroupByName(groupName, groups)
		if err != nil {
			return err
		}
		assignGroupID = group.ID.String()
	}

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

	err = s.Client.UpdateCategory(category.ID.String(), payloadName, payloadNotes, assignGroupID)
	if err != nil {
		return err
	}
	fmt.Println("Category updated with new information")
	return nil
}

func handleCategoryDelete(s *State, c *handlerContext) error {
	name, _ := c.args.pfx()

	categories, err := s.Client.GetCategories("")
	if err != nil {
		return err
	}
	category, err := findCategoryByName(name, categories)
	if err != nil {
		return err
	}

	err = s.Client.DeleteCategory(category.ID.String())
	if err != nil {
		return err
	}
	return nil
}
