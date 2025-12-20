package client

import (
	// "encoding/json"
	// "fmt"
	// "io"
	"fmt"
	"net/http"
)

type resourceNotes struct {
	Name  string `json:"name"`
	Notes string `json:"notes"`
}

// CREATE

func (c *Client) CreateBudget(name, notes string) (success bool, error error) {
	url := c.API() + "/budgets"
	payload := resourceNotes{
		Name:  name,
		Notes: notes,
	}

	resp, err := c.doRequest(http.MethodPost, url, c.LoggedInUser.JSONWebToken, payload, nil)
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
	resp, err := c.doRequest(http.MethodGet, url, c.LoggedInUser.JSONWebToken, nil, &budgets)
	if err != nil {
		return nil, err
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
