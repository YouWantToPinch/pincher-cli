package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// Get wrapper for doRequest
func (c *Client) Get(url, token string, out any) (*http.Response, error) {
	if val, ok := c.cache.Get(url); ok {
		slog.Info("retrieving requested data from cache", slog.String("url", url))
		err := json.Unmarshal(val, out)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	resp, err := c.doRequest(http.MethodGet, url, token, nil, out)
	if err == nil && out != nil {
		data, cacheErr := json.Marshal(out)
		if cacheErr != nil {
			slog.Error(fmt.Sprintf("could not cache response data: %s", cacheErr))
		} else {
			c.cache.Add(url, data)
		}
	}

	return resp, err
}

// Post wrapper for doRequest
func (c *Client) Post(url, token string, payload, out any) (*http.Response, error) {
	return c.doRequest(http.MethodPost, url, token, payload, out)
}

// Put wrapper for doRequest
func (c *Client) Put(url, token string, payload any) (*http.Response, error) {
	return c.doRequest(http.MethodPut, url, token, payload, nil)
}

// Patch wrapper for doRequest
func (c *Client) Patch(url, token string, payload any) (*http.Response, error) {
	return c.doRequest(http.MethodPatch, url, token, payload, nil)
}

// Delete wrapper for doRequest
func (c *Client) Delete(url, token string, payload any) (*http.Response, error) {
	return c.doRequest(http.MethodDelete, url, token, payload, nil)
}

func (c *Client) doRequest(method, url, token string, payload, out any) (*http.Response, error) {
	req, err := c.MakeRequest(method, url, token, payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	if out != nil {
		err := json.NewDecoder(resp.Body).Decode(out)
		if err != nil {
			return resp, err
		}
	}

	// delete existing cache for url, as the resource has been changed
	if method != http.MethodGet {
		c.cache.Delete(url)
	}

	return resp, nil
}
