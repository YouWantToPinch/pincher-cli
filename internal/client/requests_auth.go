package client

import (
	"fmt"
	"net/http"
)

// CREATE

func (c *Client) GetAccessToken(refreshToken string) (string, error) {
	url := c.API() + "/refresh"

	var accessToken string
	resp, err := c.Post(url, refreshToken, nil, &accessToken)
	if err != nil {
		return "", err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return accessToken, nil
	case http.StatusBadRequest:
		fallthrough
	case http.StatusUnauthorized:
		fallthrough
	default:
		return "", fmt.Errorf("failed to get new access token")
	}
}

// DELETE

func (c *Client) RevokeRefreshToken(refreshToken string) error {
	url := c.API() + "/revoke"

	resp, err := c.Post(url, "", nil, nil)
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
