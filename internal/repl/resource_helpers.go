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
	return nil, fmt.Errorf("no budgets found in cache with provided name")
}

func findAccountByName(name string, accounts []client.Account) (*client.Account, error) {
	for _, account := range accounts {
		if name == account.Name {
			return &account, nil
		}
	}
	return nil, fmt.Errorf("no accounts found in cache with provided name")
}
