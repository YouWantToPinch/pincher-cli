package client

import (
	"fmt"
	"log/slog"
	"net/http"
)

type logTxnSchema struct {
	AccountName         string `json:"account_name"`
	TransferAccountName string `json:"transfer_account_name"`
	// format: 2006-01-02 (YYYY-MM-DD)
	TransactionDate string           `json:"transaction_date"`
	PayeeName       string           `json:"payee_name"`
	Notes           string           `json:"notes"`
	Cleared         bool             `json:"is_cleared"`
	Amounts         map[string]int64 `json:"amounts"`
}

// CREATE

func (c *Client) LogTxn(accountName, transferAccountName, transactionDate, payeeName, notes string, isCleared bool, amounts map[string]int64) (success bool, error error) {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/transactions"
	payload := logTxnSchema{
		AccountName:         accountName,
		TransferAccountName: transferAccountName,
		TransactionDate:     transactionDate,
		PayeeName:           payeeName,
		Notes:               notes,
		Amounts:             amounts,
		Cleared:             isCleared,
	}

	resp, err := c.Post(url, c.token, payload, nil)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		if transferAccountName != "" {
			fmt.Printf("New transfer logged to account: %s\n", accountName)
		} else {
			fmt.Printf("New transaction logged to account: %s\n", accountName)
		}
		return true, nil
	default:
		return false, fmt.Errorf("failed to log transaction")
	}
}

// WARN:
// As of right now, transactions are returned IN FULL (unless query parameters are provided).
// Where this table necessarily grows rather large in size, pagination will be a must
// as soon as it is made available as a feature in the API.

func (c *Client) GetTxns(urlQuery string) ([]TransactionDetail, error) {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/transactions/details" + urlQuery
	slog.Info(url)

	type txnsContainer struct {
		Transactions []TransactionDetail `json:"transactions"`
	}

	var txns txnsContainer
	resp, cached, err := c.Get(url, c.token, &txns)
	if err != nil {
		return nil, err
	} else if cached {
		return txns.Transactions, nil
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return txns.Transactions, nil
	case http.StatusNotFound:
		return nil, fmt.Errorf("resource not found")
	default:
		return nil, fmt.Errorf("failed to retrieve budget transactions")
	}
}

// TODO:
// Add transaction updates and deletes.
//
// These are a slightly different beast, as transactions cannot be
// identified by name.
//
// The implementation is going to require something of a list view
// or other interactive approach as might be provided by Charm's 'bubbletea,'
// a courtesy that also ought be extended to the transaction LOGGING
// as well, for a better user experience.
//
// One idea might be to have a bubbletea model that renders a list of
// transactions, and when the limit is hit near the bottom, a message is
// sent outside the model to another goroutine that requests the next
// (LIMIT) amount of transactions, then adds it the the list.
