package pinchergo

import (
	"net/http"
)

func (c *Client) BudgetTransactionCreate(bID string, data BudgetTransactionCreateData) error {
	endpoint := EndpointBudgetTransactions(bID)
	var txn *Transaction
	err := c.Request(http.MethodPost, endpoint, data, &txn)
	c.Cache.addTxn(endpoint, bID, txn)
	if err != nil {
		// retrieve & cache the detailed version as well
		_, _ = c.BudgetTransactionDetails(bID, txn.ID.String())
	}
	return err
}

type transactionContainer struct {
	Transactions []*Transaction `json:"data"`
}

func (c *Client) BudgetTransaction(bID, tID string) (txn *Transaction, err error) {
	endpoint := EndpointBudgetGroup(bID, tID)
	err = c.Request(http.MethodGet, endpoint, nil, &txn)
	c.Cache.addTxn(endpoint, bID, txn)
	return txn, err
}

func (c *Client) BudgetTransactions(bID, urlQuery string) (transactions []*Transaction, err error) {
	endpoint := EndpointBudgetTransactions(bID) + urlQuery
	var container transactionContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	c.Cache.addTxns(endpoint, bID, container.Transactions)
	return container.Transactions, err
}

type transactionDetailContainer struct {
	Transactions []*TransactionDetail `json:"data"`
}

func (c *Client) BudgetTransactionDetails(bID, tID string) (txn *TransactionDetail, err error) {
	endpoint := EndpointBudgetTransactionDetails(bID, tID)
	err = c.Request(http.MethodGet, endpoint, nil, &txn)
	c.Cache.addTxnDetails(endpoint, bID, txn)
	return txn, err
}

func (c *Client) BudgetTransactionsDetails(bID, urlQuery string) (transactions []*TransactionDetail, err error) {
	endpoint := EndpointBudgetTransactionsDetails(bID) + urlQuery
	var container transactionDetailContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	c.Cache.addTxnsDetails(endpoint, bID, container.Transactions)
	return container.Transactions, err
}

func (c *Client) BudgetTransactionUpdate(bID, tID string, data BudgetTransactionUpdateData) error {
	endpoint := EndpointBudgetTransaction(bID, tID)
	err := c.Request(http.MethodPut, endpoint, data, nil)
	if err == nil {
		_, _ = c.BudgetTransaction(bID, tID)
		_, _ = c.BudgetTransactionDetails(bID, tID)
	}
	return err
}

func (c *Client) BudgetTransactionDelete(bID, tID string) error {
	endpoint := EndpointBudgetTransaction(bID, tID)
	err := c.Request(http.MethodDelete, endpoint, nil, nil)
	c.Cache.deleteTxn(bID, tID)
	c.Cache.deleteTxnsDetails(bID, tID)
	return err
}
