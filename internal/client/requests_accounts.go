package client

import (
	"net/http"
)

// BudgetAccountCreate makes an API call to create an account with the given
// data, belonging to the budget identified by the given budget ID.
//
// Any non-nil result is saved to the client's internal cache, assuming the
// cache is enabled.
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

// BudgetAccount makes an API call to get an account by ID
// from the budget identified by the given budget ID.
//
// Any non-nil result is saved to the client's internal cache, assuming the
// cache is enabled.
func (c *Client) BudgetAccount(bID, aID string) (account *Account, err error) {
	endpoint := EndpointBudgetAccount(bID, aID)
	err = c.Request(http.MethodGet, endpoint, nil, &account)
	c.Cache.addAccount(endpoint, bID, account)
	return account, err
}

// BudgetAccounts makes an API call to get accounts belonging to
// the budget identified by the given budget ID.
//
// Any non-nil result is saved to the client's internal cache, assuming the
// cache is enabled.
// The urlQuery value is appended to the endpoint string to be stored within the
// metadata of the cache entry as the Destination URL.
func (c *Client) BudgetAccounts(bID, urlQuery string) (accounts []*Account, err error) {
	endpoint := EndpointBudgetAccounts(bID) + urlQuery
	var container accountContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	c.Cache.addAccounts(endpoint, bID, container.Accounts)
	return container.Accounts, err
}

// BudgetAccountUpdate makes an API call to update an account belonging to
// the budget identified by the given budget ID with the given data.
//
// Any non-nil result is saved to the client's internal cache, assuming the
// cache is enabled. With no error returned, a separate request is made to
// fetch the account, so as to keep the cache in sync.
func (c *Client) BudgetAccountUpdate(bID, aID string, data BudgetAccountUpdateData) error {
	endpoint := EndpointBudgetAccount(bID, aID)
	err := c.Request(http.MethodPut, endpoint, data, nil)
	if err == nil {
		_, _ = c.BudgetAccount(bID, aID)
	}

	return err
}

// BudgetAccountRestore makes an API call to update an account by ID belonging to
// the budget identified by the given budget ID. The update, made via a
// PATCH request, merely alters the account's deletion status from soft-deleted
// to non-deleted.
func (c *Client) BudgetAccountRestore(bID, aID string) error {
	endpoint := EndpointBudgetAccount(bID, aID)
	err := c.Request(http.MethodPatch, endpoint, nil, nil)
	if err == nil {
		_, _ = c.BudgetAccount(bID, aID)
	}
	return err
}

// BudgetAccountDelete makes an API call to delete an account by ID belonging to
// the budget identified by the given budget ID.
//
// Any existing copy of the item in the client's internal cache is deleted,
// assuming that the payload specified hard-deltion. Otherwise, a separate
// request is made to fetch the account, so as to keep the cache in sync.
func (c *Client) BudgetAccountDelete(bID, aID string, data BudgetAccountDeleteData) error {
	endpoint := EndpointBudgetAccount(bID, aID)
	err := c.Request(http.MethodDelete, endpoint, data, nil)
	if data.DeleteHard {
		c.Cache.deleteAccount(bID, aID)
	} else {
		if err == nil {
			_, _ = c.BudgetAccount(bID, aID)
		}
	}
	return err
}
