package cli

import (
	"fmt"

	pgo "github.com/YouWantToPinch/pincher-sdk-go/pinchergo"
)

func findBudgetByName(name string, budgets []*pgo.Budget) (*pgo.Budget, error) {
	for i := range len(budgets) {
		if name == budgets[i].Name {
			return budgets[i], nil
		}
	}
	return nil, fmt.Errorf("no budgets found with provided name '%s'", name)
}

func findAccountByName(name string, accounts []*pgo.Account) (*pgo.Account, error) {
	for i := range len(accounts) {
		if name == accounts[i].Name {
			return accounts[i], nil
		}
	}
	return nil, fmt.Errorf("no accounts found with provided name '%s'", name)
}

func findGroupByName(name string, groups []*pgo.Group) (*pgo.Group, error) {
	for i := range len(groups) {
		if name == groups[i].Name {
			return groups[i], nil
		}
	}
	return nil, fmt.Errorf("no groups found with provided name '%s'", name)
}

func findCategoryByName(name string, categories []*pgo.Category) (*pgo.Category, error) {
	for i := range len(categories) {
		if name == categories[i].Name {
			return categories[i], nil
		}
	}
	return nil, fmt.Errorf("no categories found with provided name '%s'", name)
}

func findPayeeByName(name string, payees []*pgo.Payee) (*pgo.Payee, error) {
	for i := range len(payees) {
		if name == payees[i].Name {
			return payees[i], nil
		}
	}
	return nil, fmt.Errorf("no payees found with provided name '%s'", name)
}
