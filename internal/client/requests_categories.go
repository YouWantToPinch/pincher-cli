package client

import (
	"net/http"
)

func (c *Client) BudgetCategoryCreate(bID string, data BudgetCategoryCreateData) error {
	endpoint := EndpointBudgetCategories(bID)
	var category *Category
	err := c.Request(http.MethodPost, endpoint, data, &category)
	c.Cache.addCategory(bID, category)
	return err
}

func (c *Client) BudgetCategoryAssign(bID, mID string, data BudgetCategoryAssignData) error {
	endpoint := EndpointBudgetMonthCategories(bID, mID)
	err := c.Request(http.MethodPost, endpoint, data, nil)
	return err
}

type categoryContainer struct {
	Categories []*Category `json:"data"`
}

func (c *Client) BudgetCategory(bID, cID string) (category *Category, err error) {
	endpoint := EndpointBudgetCategory(bID, cID)
	err = c.Request(http.MethodGet, endpoint, nil, &category)
	c.Cache.addCategory(bID, category)
	return category, err
}

func (c *Client) BudgetCategories(bID, urlQuery string) (categories []*Category, err error) {
	endpoint := EndpointBudgetCategories(bID) + urlQuery
	var container categoryContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	c.Cache.addCategories(bID, container.Categories)
	return container.Categories, err
}

type categoryReportsContainer struct {
	CategoryReports []*CategoryReport `json:"data"`
}

func (c *Client) BudgetCategoryReports(bID, mID string) (categories []*CategoryReport, err error) {
	endpoint := EndpointBudgetMonthCategories(bID, mID)
	var container categoryReportsContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	return container.CategoryReports, err
}

func (c *Client) BudgetCategoryUpdate(bID, cID string, data BudgetCategoryUpdateData) error {
	endpoint := EndpointBudgetCategory(bID, cID)
	err := c.Request(http.MethodPut, endpoint, data, nil)
	if err == nil {
		_, _ = c.BudgetCategory(bID, cID)
	}
	return err
}

func (c *Client) BudgetCategoryRestore(bID, cID string) error {
	endpoint := EndpointBudgetCategory(bID, cID)
	err := c.Request(http.MethodPatch, endpoint, nil, nil)
	return err
}

func (c *Client) BudgetCategoryDelete(bID, cID string) error {
	endpoint := EndpointBudgetCategory(bID, cID)
	err := c.Request(http.MethodDelete, endpoint, nil, nil)
	return err
}
