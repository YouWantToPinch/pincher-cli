package client

import (
	"fmt"
	"net/http"
)

type userCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Client) CreateUser(username, password string) (success bool, error error) {
	url := c.API() + "/users"
	payload := userCredentials{
		Username: username,
		Password: password,
	}

	resp, err := c.Post(url, "", payload, nil)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return true, nil
	case http.StatusConflict:
		return false, fmt.Errorf("username already exists")
	default:
		return false, fmt.Errorf("failed to create user")
	}
}

func (c *Client) LoginUser(username, password string) (*UserInfo, error) {
	url := c.API() + "/login"
	payload := userCredentials{
		Username: username,
		Password: password,
	}

	var user UserInfo
	resp, err := c.Post(url, "", payload, &user)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return &user, nil
	default:
		return nil, fmt.Errorf("failed to log in as user")
	}
}

func (c *Client) UpdateUser(username, password string) error {
	url := c.API() + "/users"

	type rqSchema struct {
		Password string `json:"password"`
		Username string `json:"username"`
	}

	payload := rqSchema{
		Password: password,
		Username: username,
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
		return fmt.Errorf("failed to update user")
	}
}

func (c *Client) DeleteUser(username, password string) error {
	url := c.API() + "/users"
	payload := userCredentials{
		Username: username,
		Password: password,
	}

	resp, err := c.Delete(url, c.LoggedInUser.Token, payload)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	default:
		return fmt.Errorf("failed to delete user")
	}
}
