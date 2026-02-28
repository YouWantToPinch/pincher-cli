package client

import (
	"net/http"
)

func (c *Client) BudgetTransactionCreate(bID string, data BudgetTransactionCreateData) error {
	endpoint := EndpointBudgetTransactions(bID)
	err := c.Request(http.MethodPost, endpoint, data, nil)
	return err
}

type transactionContainer struct {
	Transactions []*Transaction `json:"data"`
}

func (c *Client) BudgetTransactions(bID, urlQuery string) (transactions []*Transaction, err error) {
	endpoint := EndpointBudgetTransactions(bID) + urlQuery
	var container transactionContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	return container.Transactions, err
}

type transactionDetailContainer struct {
	Transactions []*TransactionDetail `json:"data"`
}

func (c *Client) BudgetTransactionsDetails(bID, urlQuery string) (transactions []*TransactionDetail, err error) {
	endpoint := EndpointBudgetTransactionsDetails(bID) + urlQuery
	var container transactionDetailContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	return container.Transactions, err
}

func (c *Client) BudgetTransactionUpdate(bID, tID string, data BudgetTransactionUpdateData) error {
	endpoint := EndpointBudgetTransaction(bID, tID)
	err := c.Request(http.MethodPut, endpoint, data, nil)
	return err
}

func (c *Client) BudgetTransactionDelete(bID, tID string) error {
	endpoint := EndpointBudgetTransaction(bID, tID)
	err := c.Request(http.MethodDelete, endpoint, nil, nil)
	return err
}
