package client

import (
	"net/http"
)

// GetServerReady reports back with a 200 Status Code
func (c *Client) GetServerReady() (bool, error) {
	url := c.API() + "/healthz"

	resp, err := c.Get(url, "", nil)
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	default:
		return false, nil
	}
}
