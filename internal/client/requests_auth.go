package client

import (
	"fmt"
	"net/http"
)

func (c *Client) GetAccessToken(refreshToken string) (string, error) {
	url := c.API() + "/refresh"

	type rspSchema struct {
		NewAccessToken string `json:"token"`
	}

	var accessToken rspSchema
	resp, err := c.Post(url, refreshToken, nil, &accessToken)
	if err != nil {
		return "", err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return accessToken.NewAccessToken, nil
	case http.StatusBadRequest:
		fallthrough
	case http.StatusUnauthorized:
		fallthrough
	default:
		return "", fmt.Errorf("failed to get new access token")
	}
}

func (c *Client) RevokeRefreshToken(refreshToken string) error {
	url := c.API() + "/revoke"

	resp, err := c.Post(url, refreshToken, nil, nil)
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
