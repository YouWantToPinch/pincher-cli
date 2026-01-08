package repl

import (
	"fmt"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
)

func findBudgetByName(name string, budgets []client.Budget) (*client.Budget, error) {
	for i := 0; i < len(budgets); i++ {
		if name == budgets[i].Name {
			return &budgets[i], nil
		}
	}
	return nil, fmt.Errorf("no budgets found with provided name")
}

func findAccountByName(name string, accounts []client.Account) (*client.Account, error) {
	for i := 0; i < len(accounts); i++ {
		if name == accounts[i].Name {
			return &accounts[i], nil
		}
	}
	return nil, fmt.Errorf("no accounts found with provided name")
}

func findGroupByName(name string, groups []client.Group) (*client.Group, error) {
	for i := 0; i < len(groups); i++ {
		if name == groups[i].Name {
			return &groups[i], nil
		}
	}
	return nil, fmt.Errorf("no groups found with provided name")
}

func findCategoryByName(name string, categories []client.Category) (*client.Category, error) {
	for i := 0; i < len(categories); i++ {
		if name == categories[i].Name {
			return &categories[i], nil
		}
	}
	return nil, fmt.Errorf("no categories found with provided name")
}
