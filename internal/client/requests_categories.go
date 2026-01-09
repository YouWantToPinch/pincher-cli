package client

import (
	"fmt"
	"net/http"
)

func (c *Client) CreateCategory(name, notes, groupID string) (success bool, error error) {
	type rqSchema struct {
		Meta
		GroupID string `json:"group_id"`
	}

	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/categories"
	payload := rqSchema{
		Meta: Meta{
			Name:  name,
			Notes: notes,
		},
		GroupID: groupID,
	}

	resp, err := c.Post(url, c.token, payload, nil)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return true, nil
	default:
		return false, fmt.Errorf("failed to create category")
	}
}

func (c *Client) GetCategories(urlQuery string) ([]Category, error) {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/categories" + urlQuery

	type categoryContainer struct {
		Categories []Category `json:"categories"`
	}

	var categories categoryContainer
	resp, err := c.Get(url, c.token, &categories)
	if err != nil {
		return nil, err
	} else if resp == nil {
		return categories.Categories, nil
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return categories.Categories, nil
	case http.StatusNotFound:
		return nil, fmt.Errorf("resource not found")
	default:
		return nil, fmt.Errorf("failed to retrieve budget categories")
	}
}

func (c *Client) UpdateCategory(categoryID, name, notes, groupID string) error {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/categories/" + categoryID

	type rqSchema struct {
		Meta
		GroupID string `json:"group_id"`
	}

	payload := rqSchema{
		Meta: Meta{
			Name:  name,
			Notes: notes,
		},
		GroupID: groupID,
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
		return fmt.Errorf("failed to retrieve budget categories")
	}
}

func (c *Client) DeleteCategory(categoryID string) error {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/categories/" + categoryID

	resp, err := c.Delete(url, c.token, nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Category deleted. It may be restored, or permanently deleted.")
		return nil
	case http.StatusNoContent:
		fmt.Println("Category deleted. It cannot be restored.")
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("resource not found")
	case http.StatusBadRequest:
		return fmt.Errorf("bad request (has.Category been soft-deleted first?)")
	default:
		return fmt.Errorf("failed to retrieve budget categories")
	}
}
