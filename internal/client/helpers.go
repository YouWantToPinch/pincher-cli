package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
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
	switch method {
	case http.MethodPost:
		c.cache.Delete(url)
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		path, err := getURLResourcePath(url)
		if err != nil {
			slog.Error("could not delete key from url entry cache", slog.String("key", url))
		}
		c.cache.Delete(path)
	}

	return resp, nil
}

// getUrlResourcePath trims ONE trailing instance of "/..."
// from a string with at least one "/".
// It is made for the purpose of trimming a trailing resource ID
// from the end of a URL in order to get its parent path.
// This is useful if you want to clear a cache of this resource
// type so as to avoid cache inconsistency.
func getURLResourcePath(url string) (string, error) {
	// Remove the last segment of the path
	segments := strings.Split(url, "/")
	if len(segments) > 1 {
		url = strings.Join(segments[:len(segments)-1], "/")
		return url, nil
	} else {
		return "", fmt.Errorf("input is not a url")
	}
}
