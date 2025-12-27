package repl

import (
	"fmt"
	"sort"
	"strings"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
)

func handlerGroup(s *State, c *handlerContext) error {
	if val, ok := c.ctxValues["action"]; ok {
		switch val {
		case "add":
			return handleGroupAdd(s, c)
		case "list":
			return handleGroupList(s, c)
		case "update":
			return handleGroupUpdate(s, c)
		case "delete":
			return handleGroupDelete(s, c)
		default:
			return fmt.Errorf("action not implemented")
		}
	} else {
		return fmt.Errorf("action was not saved to context")
	}
}

func handleGroupAdd(s *State, c *handlerContext) error {
	name, _ := c.args.pfx()

	c.args.trackOptArgs(&c.cmd, "notes")
	notes, _ := c.args.pfx()

	budgetCreated, err := s.Client.CreateGroup(name, notes)
	if err != nil {
		return err
	}
	if budgetCreated {
		fmt.Println("Group " + name + " successfully created as user: " + s.Client.LoggedInUser.Username + ".")
		fmt.Println("See it with: `group list`")
		return nil
	} else {
		return fmt.Errorf("budget could not be created")
	}
}

func handleGroupList(s *State, c *handlerContext) error {
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

	groups, err := s.Client.GetGroups(includeQuery)
	if err != nil {
		return err
	}
	if len(groups) == 0 {
		fmt.Printf("No groups found belonging to budget %s. \n", s.Client.ViewedBudget.Name)
		return nil
	}
	fmt.Printf("Groups under budget %s: \n", s.Client.ViewedBudget.Name)
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	const uuidLength = 36
	maxLenName := MaxOfStrings(ExtractStrings(groups, func(b client.Group) string { return b.Name }))
	maxLenNotes := MaxOfStrings(ExtractStrings(groups, func(b client.Group) string { return b.Notes }))
	fmt.Printf("  %-*s | %-*s | %s\n", maxLenName, "NAME", uuidLength, "ID", "NOTES")
	fmt.Printf("  %s-+-%s-+-%s\n", nDashes(maxLenName), nDashes(uuidLength), nDashes(maxLenNotes))
	for _, group := range groups {
		fmt.Printf("  %-*s  %-*s   %s\n", maxLenName, group.Name, uuidLength, group.ID, group.Notes)
	}

	return nil
}

func handleGroupUpdate(s *State, c *handlerContext) error {
	groupName, _ := c.args.pfx()

	groups, err := s.Client.GetGroups("")
	if err != nil {
		return err
	}
	group, err := findGroupByName(groupName, groups)
	if err != nil {
		return err
	}

	c.args.trackOptArgs(&c.cmd, "name")
	payloadName, err := c.args.pfx()
	if err != nil {
		payloadName = group.Name
	}
	c.args.trackOptArgs(&c.cmd, "notes")
	payloadNotes, err := c.args.pfx()
	if err != nil {
		payloadNotes = group.Notes
	}

	err = s.Client.UpdateGroup(group.ID.String(), payloadName, payloadNotes)
	if err != nil {
		return err
	}
	fmt.Println("Group updated with new information")
	return nil
}

func handleGroupDelete(s *State, c *handlerContext) error {
	name, _ := c.args.pfx()
	c.args.trackOptArgs(&c.cmd, "hard")
	flagDeleteHard, _ := c.args.pfx()

	groups, err := s.Client.GetGroups("")
	if err != nil {
		return err
	}
	group, err := findGroupByName(name, groups)
	if err != nil {
		return err
	}

	err = s.Client.DeleteGroup(group.ID.String(), name, flagDeleteHard)
	if err != nil {
		return err
	}
	fmt.Println("Group deleted.")
	return nil
}
