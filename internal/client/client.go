package client

import (
	"bytes"
	"encoding/json"
	pcache "github.com/YouWantToPinch/pincher-cli/internal/pinchercache"
	"io"
	"net/http"
	"time"
)

type userInfo struct {
	JSONWebToken string
	Username     string
}

type Client struct {
	cache        pcache.Cache
	httpClient   http.Client
	LoggedInUser userInfo
	BaseUrl      string
}

func (c *Client) API() string {
	return c.BaseUrl + "/api"
}

func NewClient(timeout, cacheInterval time.Duration, baseUrl string) Client {
	var url string
	if baseUrl == "" {
		url = defaultBaseUrl
	} else {
		url = baseUrl
	}
	return Client{
		cache: pcache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
		BaseUrl: url,
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
	// TODO: Handle error, log to a local directory
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return req, nil
}
