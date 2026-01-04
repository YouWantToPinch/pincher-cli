package repl

import (
	"fmt"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
)

func findBudgetByName(name string, budgets []client.Budget) (*client.Budget, error) {
	for _, budget := range budgets {
		if name == budget.Name {
			return &budget, nil
		}
	}
	return nil, fmt.Errorf("no budgets found with provided name")
}

func findAccountByName(name string, accounts []client.Account) (*client.Account, error) {
	for _, account := range accounts {
		if name == account.Name {
			return &account, nil
		}
	}
	return nil, fmt.Errorf("no accounts found with provided name")
}

func findGroupByName(name string, groups []client.Group) (*client.Group, error) {
	for _, group := range groups {
		if name == group.Name {
			return &group, nil
		}
	}
	return nil, fmt.Errorf("no groups found with provided name")
}

func findCategoryByName(name string, categories []client.Category) (*client.Category, error) {
	for _, category := range categories {
		if name == category.Name {
			return &category, nil
		}
	}
	return nil, fmt.Errorf("no categories found with provided name")
}
