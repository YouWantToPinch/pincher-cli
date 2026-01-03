package client

import (
	"fmt"
	"net/http"
)

// CREATE

func (c *Client) CreateBudget(name, notes string) (success bool, error error) {
	url := c.API() + "/budgets"

	payload := Meta{
		Name:  name,
		Notes: notes,
	}

	resp, err := c.Post(url, c.LoggedInUser.Token, payload, nil)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return true, nil
	default:
		return false, fmt.Errorf("failed to create budget")
	}
}

func (c *Client) GetBudgets(urlQuery string) ([]Budget, error) {
	url := c.API() + "/budgets" + urlQuery

	type budgetContainer struct {
		Budgets []Budget `json:"budgets"`
	}

	var budgets budgetContainer
	resp, err := c.Get(url, c.LoggedInUser.Token, &budgets)
	if err != nil {
		return nil, err
	} else if resp == nil {
		return budgets.Budgets, nil
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return budgets.Budgets, nil
	case http.StatusNotFound:
		return nil, fmt.Errorf("resource not found")
	default:
		return nil, fmt.Errorf("failed to retrieve user budgets")
	}
}
