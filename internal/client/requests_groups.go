package client

import (
	"net/http"
)

func (c *Client) BudgetGroupCreate(bID string, data BudgetGroupCreateData) error {
	endpoint := EndpointBudgetGroups(bID)
	var group *Group
	err := c.Request(http.MethodPost, endpoint, data, &group)
	c.Cache.addGroup(bID, group)
	return err
}

type groupContainer struct {
	Groups []*Group `json:"data"`
}

func (c *Client) BudgetGroup(bID, gID string) (group *Group, err error) {
	endpoint := EndpointBudgetGroup(bID, gID)
	err = c.Request(http.MethodGet, endpoint, nil, &group)
	c.Cache.addGroup(bID, group)
	return group, err
}

func (c *Client) BudgetGroups(bID, urlQuery string) (groups []*Group, err error) {
	endpoint := EndpointBudgetGroups(bID) + urlQuery
	var container groupContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	c.Cache.addGroups(bID, container.Groups)
	return container.Groups, err
}

func (c *Client) BudgetGroupUpdate(bID, gID string, data BudgetGroupUpdateData) error {
	endpoint := EndpointBudgetGroup(bID, gID)
	err := c.Request(http.MethodPut, endpoint, data, nil)
	if err == nil {
		_, _ = c.BudgetGroup(bID, gID)
	}
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
