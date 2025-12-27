package client

import (
	"fmt"
	"net/http"
)

// CREATE

func (c *Client) CreateAccount(name, notes, accountType string) (success bool, error error) {
	type rqSchema struct {
		Meta
		AccountType string `json:"account_type"`
	}

	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/accounts"
	payload := rqSchema{
		Meta: Meta{
			Name:  name,
			Notes: notes,
		},
		AccountType: accountType,
	}

	resp, err := c.Post(url, c.LoggedInUser.JSONWebToken, payload, nil)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return true, nil
	default:
		return false, fmt.Errorf("failed to create account")
	}
}

func (c *Client) GetAccounts(urlQuery string) ([]Account, error) {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/accounts" + urlQuery

	type accountContainer struct {
		Accounts []Account `json:"accounts"`
	}

	var accounts accountContainer
	resp, err := c.Get(url, c.LoggedInUser.JSONWebToken, &accounts)
	if err != nil {
		return nil, err
	} else if resp == nil {
		return accounts.Accounts, nil
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return accounts.Accounts, nil
	case http.StatusNotFound:
		return nil, fmt.Errorf("resource not found")
	default:
		return nil, fmt.Errorf("failed to retrieve budget accounts")
	}
}

func (c *Client) UpdateAccount(accountID, name, notes, accountType string) error {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/accounts/" + accountID

	type rqSchema struct {
		Meta
		AccountType string `json:"account_id"`
	}

	payload := rqSchema{
		Meta: Meta{
			Name:  name,
			Notes: notes,
		},
		AccountType: accountType,
	}

	resp, err := c.Put(url, c.LoggedInUser.JSONWebToken, payload)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("resource not found")
	default:
		return fmt.Errorf("failed to retrieve budget accounts")
	}
}

func (c *Client) RestoreAccount(accountID, name string) error {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/accounts/" + accountID

	resp, err := c.Patch(url, c.LoggedInUser.JSONWebToken, nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("resource not found")
	default:
		return fmt.Errorf("failed to retrieve budget accounts")
	}
}

func (c *Client) DeleteAccount(accountID, name, deleteHard string) error {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/accounts/" + accountID

	type rqSchema struct {
		Name       string `json:"name"`
		DeleteHard bool   `json:"delete_hard"`
	}

	payload := rqSchema{
		Name:       name,
		DeleteHard: deleteHard == "SET",
	}

	resp, err := c.Delete(url, c.LoggedInUser.JSONWebToken, payload)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("resource not found")
	case http.StatusBadRequest:
		return fmt.Errorf("bad request (has account been soft-deleted first?)")
	default:
		return fmt.Errorf("failed to retrieve budget accounts")
	}
}
