package client

import (
	"fmt"
	"net/http"
)

func (c *Client) CreateBudget(name, notes string) (success bool, error error) {
	url := c.API() + "/budgets"

	payload := Meta{
		Name:  name,
		Notes: notes,
	}

	resp, err := c.Post(url, c.token, payload, nil)
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
	resp, err := c.Get(url, c.token, &budgets)
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
		return nil, fmt.Errorf("failed to retrieve budgets")
	}
}

func (c *Client) UpdateBudget(budgetID, name, notes string) error {
	url := c.API() + "/budgets/" + budgetID

	type rqSchema struct {
		Meta
	}

	payload := rqSchema{
		Meta: Meta{
			Name:  name,
			Notes: notes,
		},
	}

	resp, err := c.Put(url, c.token, payload)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("resource not found")
	default:
		return fmt.Errorf("failed to update budget")
	}
}

func (c *Client) DeleteBudget(budgetID string) error {
	url := c.API() + "/budgets/" + budgetID

	resp, err := c.Delete(url, c.token, nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("resource not found")
	default:
		return fmt.Errorf("failed to delete budgets")
	}
}
