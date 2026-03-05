package client

import (
	"net/http"
)

func (c *Client) BudgetAccountCreate(bID string, data BudgetAccountCreateData) error {
	endpoint := EndpointBudgetAccounts(bID)
	var account *Account
	err := c.Request(http.MethodPost, endpoint, data, &account)
	c.Cache.addAccount(endpoint, bID, account)
	return err
}

type accountContainer struct {
	Accounts []*Account `json:"data"`
}

func (c *Client) BudgetAccount(bID, aID string) (account *Account, err error) {
	endpoint := EndpointBudgetAccount(bID, aID)
	err = c.Request(http.MethodGet, endpoint, nil, &account)
	c.Cache.addAccount(endpoint, bID, account)
	return account, err
}

func (c *Client) BudgetAccounts(bID, urlQuery string) (accounts []*Account, err error) {
	endpoint := EndpointBudgetAccounts(bID) + urlQuery
	var container accountContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	c.Cache.addAccounts(endpoint, bID, container.Accounts)
	return container.Accounts, err
}

func (c *Client) BudgetAccountUpdate(bID, aID string, data BudgetAccountUpdateData) error {
	endpoint := EndpointBudgetAccount(bID, aID)
	err := c.Request(http.MethodPut, endpoint, data, nil)
	if err == nil {
		_, _ = c.BudgetAccount(bID, aID)
	}

	return err
}

func (c *Client) BudgetAccountRestore(bID, aID string) error {
	endpoint := EndpointBudgetAccount(bID, aID)
	err := c.Request(http.MethodPatch, endpoint, nil, nil)
	return err
}

func (c *Client) BudgetAccountDelete(bID, aID string, data BudgetAccountDeleteData) error {
	endpoint := EndpointBudgetAccount(bID, aID)
	err := c.Request(http.MethodDelete, endpoint, data, nil)
	c.Cache.deleteAccount(bID, aID)
	return err
}
