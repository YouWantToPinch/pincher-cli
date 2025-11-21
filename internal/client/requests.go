package client

import (
	// "encoding/json"
	// "fmt"
	// "io"
	"net/http"
)

func (c *Client) GetServerReady() (bool, error) {
	url := c.API() + "/healthz"

	// make a request
	req, err := c.MakeRequest(http.MethodGet, url, "", nil)
	if err != nil {
		return false, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, nil
}
