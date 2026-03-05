package pinchergo

import (
	"net/http"
)

// GetServerReady reports back with a 200 Status Code
func (c *Client) GetServerReady() (bool, error) {
	endpoint := EndpointServerReadiness()
	err := c.Request(http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}
