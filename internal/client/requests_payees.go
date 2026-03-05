package client

import (
	"net/http"
)

func (c *Client) BudgetPayeeCreate(bID string, data BudgetPayeeCreateData) error {
	endpoint := EndpointBudgetPayees(bID)
	var payee *Payee
	err := c.Request(http.MethodPost, endpoint, data, &payee)
	c.Cache.addPayee(bID, payee)
	return err
}

type payeeContainer struct {
	Payees []*Payee `json:"data"`
}

func (c *Client) BudgetPayee(bID, pID string) (payee *Payee, err error) {
	endpoint := EndpointBudgetPayee(bID, pID)
	err = c.Request(http.MethodGet, endpoint, nil, &payee)
	c.Cache.addPayee(bID, payee)
	return payee, err
}

func (c *Client) BudgetPayees(bID, urlQuery string) (Payees []*Payee, err error) {
	endpoint := EndpointBudgetPayees(bID) + urlQuery
	var container payeeContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	c.Cache.addPayees(bID, container.Payees)
	return container.Payees, err
}

func (c *Client) BudgetPayeeUpdate(bID, pID string, data BudgetPayeeUpdateData) error {
	endpoint := EndpointBudgetPayee(bID, pID)
	err := c.Request(http.MethodPut, endpoint, data, nil)
	if err == nil {
		_, _ = c.BudgetPayee(bID, pID)
	}
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
