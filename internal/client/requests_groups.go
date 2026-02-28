package client

import (
	"net/http"
)

func (c *Client) BudgetGroupCreate(bID string, data BudgetGroupCreateData) error {
	endpoint := EndpointBudgetGroups(bID)
	err := c.Request(http.MethodPost, endpoint, data, nil)
	return err
}

type groupContainer struct {
	Groups []*Group `json:"data"`
}

func (c *Client) BudgetGroups(bID, urlQuery string) (groups []*Group, err error) {
	endpoint := EndpointBudgetGroups(bID) + urlQuery
	var container groupContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	return container.Groups, err
}

func (c *Client) BudgetGroupUpdate(bID, gID string, data BudgetGroupUpdateData) error {
	endpoint := EndpointBudgetGroup(bID, gID)
	err := c.Request(http.MethodPut, endpoint, data, nil)
	return err
}

func (c *Client) BudgetGroupRestore(bID, gID string) error {
	endpoint := EndpointBudgetGroup(bID, gID)
	err := c.Request(http.MethodPatch, endpoint, nil, nil)
	return err
}

func (c *Client) BudgetGroupDelete(bID, gID string) error {
	endpoint := EndpointBudgetGroup(bID, gID)
	err := c.Request(http.MethodDelete, endpoint, nil, nil)
	return err
}
