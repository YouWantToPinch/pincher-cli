package client

import (
	"fmt"
	"net/http"
)

func (c *Client) CreatePayee(name, notes string) (success bool, error error) {
	type rqSchema struct {
		Meta
	}

	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/payees"
	payload := rqSchema{
		Meta: Meta{
			Name:  name,
			Notes: notes,
		},
	}

	resp, err := c.Post(url, c.token, payload, nil)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return true, nil
	default:
		return false, fmt.Errorf("failed to create payee")
	}
}

func (c *Client) GetPayees(urlQuery string) ([]Payee, error) {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/payees" + urlQuery

	type payeeContainer struct {
		Payees []Payee `json:"payees"`
	}

	var payees payeeContainer
	resp, cached, err := c.Get(url, c.token, &payees)
	if err != nil {
		return nil, err
	} else if cached {
		return payees.Payees, nil
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return payees.Payees, nil
	case http.StatusNotFound:
		return nil, fmt.Errorf("resource not found")
	default:
		return nil, fmt.Errorf("failed to retrieve budget payees")
	}
}

func (c *Client) UpdatePayee(payeeID, name, notes string) error {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/payees/" + payeeID

	type rqSchema struct {
		Meta
	}

	payload := rqSchema{
		Meta: Meta{
			Name:  name,
			Notes: notes,
		},
	}

	resp, err := c.Put(url, c.token, payload)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("resource not found")
	default:
		return fmt.Errorf("failed to retrieve budget payees")
	}
}

func (c *Client) DeletePayee(payeeID, newPayeeName string) error {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/payees/" + payeeID

	type rqSchema struct {
		NewPayeeName string `json:"new_payee_name"`
	}

	payload := rqSchema{
		NewPayeeName: newPayeeName,
	}

	resp, err := c.Delete(url, c.token, payload)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("resource not found")
	default:
		return fmt.Errorf("failed to retrieve budget payees")
	}
}
