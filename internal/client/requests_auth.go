package client

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (c *Client) UserTokenRefreshWithUser() (user *User, err error) {
	if c.RefreshToken == "" {
		slog.Warn("Client directed to get new access token with user, but Refresh Token was empty.")
		return nil, fmt.Errorf("refresh token is empty")
	} else {
		fmt.Println("RefreshToken: " + c.RefreshToken)
	}

	endpoint := EndpointRefresh() + "?with-user"

	type rspSchema struct {
		User
		NewAccessToken string `json:"token"`
	}

	var resp rspSchema
	err = c.RequestWithToken(&c.RefreshToken, http.MethodPost, endpoint, nil, &resp)
	if err != nil {
		return nil, err
	}

	c.token = resp.NewAccessToken
	return &resp.User, nil
}

func (c *Client) UserTokenRefresh() error {
	if c.RefreshToken == "" {
		slog.Warn("Client directed to get new access token, but Refresh Token was empty.")
		return fmt.Errorf("refresh token is empty")
	}

	endpoint := EndpointRefresh()

	type rspSchema struct {
		NewAccessToken string `json:"token"`
	}

	var token rspSchema

	err := c.RequestWithToken(&c.RefreshToken, http.MethodPost, endpoint, nil, &token)
	if err != nil {
		return err
	}

	c.token = token.NewAccessToken
	return nil
}

func (c *Client) UserTokenRevoke() error {
	if c.RefreshToken == "" {
		slog.Warn("Client directed to revoke active refresh token, but found none to revoke.")
		return nil
	}

	endpoint := EndpointRevoke()
	err := c.RequestWithToken(&c.RefreshToken, http.MethodPost, endpoint, nil, nil)

	// Whether or not the server has trouble revoking,
	// we can at least forget it here in the CLI as well.
	c.RefreshToken = ""

	return err
}
