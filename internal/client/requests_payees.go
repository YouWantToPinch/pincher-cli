package client

import (
	"net/http"
)

func (c *Client) BudgetPayeeCreate(bID string, data BudgetPayeeCreateData) error {
	endpoint := EndpointBudgetPayees(bID)
	err := c.Request(http.MethodPost, endpoint, data, nil)
	return err
}

type payeeContainer struct {
	Payees []*Payee `json:"data"`
}

func (c *Client) BudgetPayees(bID, urlQuery string) (Payees []*Payee, err error) {
	endpoint := EndpointBudgetPayees(bID) + urlQuery
	var container payeeContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	return container.Payees, err
}

func (c *Client) BudgetPayeeUpdate(bID, pID string, data BudgetPayeeUpdateData) error {
	endpoint := EndpointBudgetPayee(bID, pID)
	err := c.Request(http.MethodPut, endpoint, data, nil)
	return err
}

func (c *Client) BudgetPayeeRestore(bID, pID string) error {
	endpoint := EndpointBudgetPayee(bID, pID)
	err := c.Request(http.MethodPatch, endpoint, nil, nil)
	return err
}

func (c *Client) BudgetPayeeDelete(bID, pID string, data BudgetPayeeDeleteData) error {
	endpoint := EndpointBudgetPayee(bID, pID)
	err := c.Request(http.MethodDelete, endpoint, data, nil)
	return err
}
