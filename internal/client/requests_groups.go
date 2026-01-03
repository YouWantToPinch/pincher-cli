package client

import (
	"fmt"
	"net/http"
)

// CREATE

func (c *Client) CreateGroup(name, notes string) (success bool, error error) {
	type rqSchema struct {
		Meta
	}

	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/groups"
	payload := rqSchema{
		Meta: Meta{
			Name:  name,
			Notes: notes,
		},
	}

	resp, err := c.Post(url, c.LoggedInUser.Token, payload, nil)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return true, nil
	default:
		return false, fmt.Errorf("failed to create group")
	}
}

func (c *Client) GetGroups(urlQuery string) ([]Group, error) {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/groups" + urlQuery

	type groupContainer struct {
		Groups []Group `json:"groups"`
	}

	var groups groupContainer
	resp, err := c.Get(url, c.LoggedInUser.Token, &groups)
	if err != nil {
		return nil, err
	} else if resp == nil {
		return groups.Groups, nil
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return groups.Groups, nil
	case http.StatusNotFound:
		return nil, fmt.Errorf("resource not found")
	default:
		return nil, fmt.Errorf("failed to retrieve budget groups")
	}
}

func (c *Client) UpdateGroup(groupID, name, notes string) error {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/groups/" + groupID

	type rqSchema struct {
		Meta
	}

	payload := rqSchema{
		Meta: Meta{
			Name:  name,
			Notes: notes,
		},
	}

	resp, err := c.Put(url, c.LoggedInUser.Token, payload)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("resource not found")
	default:
		return fmt.Errorf("failed to retrieve budget groups")
	}
}

func (c *Client) DeleteGroup(groupID, name string) error {
	url := c.API() + "/budgets/" + c.ViewedBudget.ID.String() + "/groups/" + groupID

	resp, err := c.Delete(url, c.LoggedInUser.Token, nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("resource not found")
	default:
		return fmt.Errorf("failed to retrieve budget groups")
	}
}
