// Package client handles pincher-api calls
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	file "github.com/YouWantToPinch/pincher-cli/internal/filemgr"
	"github.com/golang-jwt/jwt/v5"
)

type Client struct {
	http.Client
	cache        Cache
	ViewedBudget Budget
	BaseURL      string
	token        string
	RefreshToken string
}

// ClearUserSession attempts to revoke the refresh token in the current session,
// and then clears the client and cache of its LoggedInUser values, before
// then saving the cache to disk.
func (c *Client) ClearUserSession() {
	revokeErr := c.RevokeRefreshToken()
	revokeResult := "success"
	if revokeErr != nil {
		revokeResult = "failure"
	}
	slog.Info("Attempted to revoke refresh token with result: " + revokeResult)
	c.cache.Clear()
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
		cache: *NewCache(cacheInterval),
		Client: http.Client{
			Timeout: timeout,
		},
		BaseURL: url,
	}
}

// LoadCacheFile looks for a file with the given name within
// the user cache directory and attempts to load it into the cache.
func (c *Client) LoadCacheFile() error {
	const errMsg = "could not load cache: "
	cachePath, err := file.GetCacheFilepath("cache.json")
	if err != nil {
		return fmt.Errorf(errMsg+"%w", err)
	}

	loadedCache, err := file.ReadJSONFromFile[Cache](cachePath)
	if err != nil {
		return fmt.Errorf(errMsg+"%w", err)
	}

	c.cache.Set(loadedCache.CachedEntries)
	return nil
}

// SaveCacheFile saves the current cache to a local file
// with the given name under the user cache directory.
func (c *Client) SaveCacheFile() error {
	const errMsg = "could not save cache: "
	cachePath, err := file.GetCacheFilepath("cache.json")
	if err != nil {
		return fmt.Errorf(errMsg+"%w", err)
	}
	err = file.WriteAsJSON(c.cache, cachePath)
	if err != nil {
		return fmt.Errorf(errMsg+"%w", err)
	}
	return nil
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

// Get wrapper for doRequest
func (c *Client) Get(url, token string, out any) (*http.Response, error) {
	if val, ok := c.cache.Get(url); ok {
		slog.Info("retrieving requested data from cache", slog.String("URL", url))
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
			c.cache.Add(url, data, false)
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
	// try to make the new request.
	req, err := c.MakeRequest(method, url, token, payload)
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	shouldRetry := !strings.Contains(url, "/refresh") && !strings.Contains(url, "/revoke") // avoid infinite loop
	for {
		// try to send the request.
		resp, err = c.Do(req)
		if err != nil {
			return nil, err
		}

		// If the server responds with a 401 (unauthorized) code,
		// we need to check whether or not it concerns an expired token.
		tokenExpired, err := checkTokenExpired(token)

		// if it isn't a 401, we don't care to retry with a new access token
		if resp.StatusCode != http.StatusUnauthorized ||
			// if no token was required for the request, we don't care
			token == "" ||
			// if the token is not expired, the 401 has nothing to do with tokens
			(!tokenExpired && err == nil) {
			break
		} else if shouldRetry {

			// Try to get a new access token if it is invalid or we got an error.
			// If that's not possible, log out the user, as their session must therefore be invalid.
			success, err := c.GetAccessToken()
			if !success || err != nil {
				// an error return here means we've been logged out
				return nil, err
			}
			// update this loop with the new access token
			token = c.token
			shouldRetry = false
		} else {
			break
		}
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
		c.cache.DeleteAllStartsWith(url)
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		path, err := getURLResourcePath(url)
		if err != nil {
			slog.Error("could not delete keys with url prefix from cache", slog.String("prefix", url))
		} else {
			// delete all probable cache of this resource
			c.cache.DeleteAllStartsWith(path)
		}
	}

	return resp, nil
}

// --------------
//  Helpers
// --------------

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

func checkTokenExpired(tokenString string) (bool, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("invalid claims")
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return false, err
	}

	return exp.Before(time.Now()), nil
}
