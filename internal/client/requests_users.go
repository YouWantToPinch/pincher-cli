package client

import (
	// "encoding/json"
	// "fmt"
	// "io"
	"fmt"
	"net/http"
)

type userCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CREATE
func (c *Client) CreateUser(username, password string) (success bool, error error) {
	url := c.API() + "/users"
	payload := userCredentials{
		Username: username,
		Password: password,
	}

	resp, err := c.doRequest(http.MethodPost, url, "", payload, nil)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return true, nil
	case http.StatusConflict:
		return false, fmt.Errorf("Username already exists.")
	default:
		return false, fmt.Errorf("Failed to create user.")
	}
}

func (c *Client) LoginUser(username, password string) (User, error) {
	url := c.API() + "/login"
	payload := userCredentials{
		Username: username,
		Password: password,
	}

	var user User
	resp, err := c.doRequest(http.MethodPost, url, "", payload, &user)
	if err != nil {
		return User{}, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return user, nil
	default:
		return User{}, fmt.Errorf("Failed to log in as user")
	}
}
