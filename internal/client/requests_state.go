package client

import (
	// "encoding/json"
	// "fmt"
	// "io"
	"net/http"
)

// Admin & State
func (c *Client) GetServerReady() (bool, error) {
	url := c.API() + "/healthz"

	resp, err := c.doRequest(http.MethodGet, url, "", nil, nil)
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
