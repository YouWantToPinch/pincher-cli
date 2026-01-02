// Package client handles pincher-api calls
package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	cache "github.com/YouWantToPinch/pincher-cli/internal/cache"
)

type userInfo struct {
	JSONWebToken string
	Username     string
	ID           string
}

type Client struct {
	cache        cache.Cache
	httpClient   http.Client
	LoggedInUser userInfo
	ViewedBudget Budget
	BaseURL      string
}

func (c *Client) API() string {
	return c.BaseURL + "/api"
}

func NewClient(timeout, cacheInterval time.Duration, baseURL string) Client {
	var url string
	if baseURL == "" {
		url = defaultBaseURL
	} else {
		url = baseURL
	}
	return Client{
		cache: cache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
		BaseURL: url,
	}
}

func (c *Client) MakeRequest(method, path, token string, body any) (*http.Request, error) {
	var buffer io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buffer = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, path, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return req, nil
}
