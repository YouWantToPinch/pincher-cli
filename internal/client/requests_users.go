package client

import (
	"net/http"
)

func (c *Client) UserCreate(data UserCreateData) error {
	endpoint := EndpointUsers()
	err := c.Request(http.MethodPost, endpoint, data, nil)
	return err
}

func (c *Client) UserLogin(data UserLoginData) (user *User, err error) {
	endpoint := EndpointLogin()

	type rspSchema struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	var login rspSchema

	err = c.Request(http.MethodPost, endpoint, data, &login)
	if err != nil {
		return nil, err
	}

	c.token = login.Token
	c.RefreshToken = login.RefreshToken
	return &login.User, err
}

func (c *Client) UserUpdate(data UserUpdateData) error {
	endpoint := EndpointUsers()
	err := c.Request(http.MethodPost, endpoint, data, nil)
	return err
}

func (c *Client) UserDelete(data UserDeleteData) error {
	endpoint := EndpointUsers()
	err := c.Request(http.MethodDelete, endpoint, data, nil)
	return err
}
