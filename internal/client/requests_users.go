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

func (c *Client) LoginUser(username, password string) (User, error) {
	url := c.API() + "/login"
	payload := userCredentials{
		Username: username,
		Password: password,
	}

	type rspSchema struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	var login rspSchema
	resp, err := c.Post(url, "", payload, &login)
	if err != nil {
		return User{}, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		c.token = login.Token
		c.RefreshToken = login.RefreshToken
		return login.User, nil
	case http.StatusUnauthorized:
		return User{}, fmt.Errorf("incorrect username or password")
	default:
		return User{}, fmt.Errorf("could not log in")
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
		return fmt.Errorf("failed to update user")
	}
}

func (c *Client) DeleteUser(username, password string) error {
	url := c.API() + "/users"
	payload := userCredentials{
		Username: username,
		Password: password,
	}

	resp, err := c.Delete(url, c.token, payload)
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
