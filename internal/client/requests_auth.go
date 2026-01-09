package client

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (c *Client) GetAccessToken() (success bool, err error) {
	url := c.API() + "/refresh"

	type rspSchema struct {
		NewAccessToken string `json:"token"`
	}

	var token rspSchema
	resp, err := c.Post(url, c.RefreshToken, nil, &token)
	if err != nil {
		return false, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		c.token = token.NewAccessToken
		return true, nil
	case http.StatusBadRequest:
		fallthrough
	case http.StatusUnauthorized:
		fallthrough
	default:
		return false, fmt.Errorf("could not get new access token")
	}
}

func (c *Client) RevokeRefreshToken() error {
	if c.RefreshToken == "" {
		slog.Warn("Client directed to revoke active refresh token, but found none to revoke.")
		return nil
	}

	url := c.API() + "/revoke"

	resp, err := c.Post(url, c.RefreshToken, nil, nil)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		fallthrough
	case http.StatusUnauthorized:
		fallthrough
	default:
		return fmt.Errorf("failed to revoke refresh token")
	}
}
