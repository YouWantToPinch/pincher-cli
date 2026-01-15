package client

import (
	"fmt"
	"net/http"
)

func (c *Client) CreateCategory(name, notes, groupName string) (success bool, error error) {
	type rqSchema struct {
		Meta
		GroupName string `json:"group_name"`
	}

	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/categories"
	payload := rqSchema{
		Meta: Meta{
			Name:  name,
			Notes: notes,
		},
		GroupName: groupName,
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

func (c *Client) AssignAmountToCategory(amount int64, toCategoryName, fromCategoryName, monthID string) error {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/months/" + monthID + "/categories"

	type rqSchema struct {
		Amount       int64  `json:"amount"`
		ToCategory   string `json:"to_category"`
		FromCategory string `json:"from_category"`
	}
	payload := &rqSchema{
		Amount:       amount,
		ToCategory:   toCategoryName,
		FromCategory: fromCategoryName,
	}

	resp, err := c.Post(url, c.token, payload, nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return fmt.Errorf("improper input")
	default:
		return fmt.Errorf("%d: failed to assign amount to category", resp.StatusCode)
	}
}

func (c *Client) GetCategoryReports(monthID string) ([]CategoryReport, error) {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/months/" + monthID + "/categories"

	type rspSchema struct {
		CategoryReports []CategoryReport `json:"category_reports"`
	}

	var rspPayload rspSchema
	resp, cached, err := c.Get(url, c.token, &rspPayload)
	if err != nil {
		return nil, err
	} else if cached {
		return rspPayload.CategoryReports, nil
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return rspPayload.CategoryReports, nil
	case http.StatusBadRequest:
		return nil, fmt.Errorf("improper input")
	default:
		return nil, fmt.Errorf("failed to assign amount to category")
	}
}

func (c *Client) GetCategories(urlQuery string) ([]Category, error) {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/categories" + urlQuery

	type rspSchema struct {
		Categories []Category `json:"categories"`
	}

	var rspPayload rspSchema
	resp, cached, err := c.Get(url, c.token, &rspPayload)
	if err != nil {
		return nil, err
	} else if cached {
		return rspPayload.Categories, nil
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return rspPayload.Categories, nil
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
