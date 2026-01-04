package client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Get wrapper for doRequest
func (c *Client) Get(url, token string, out any) (*http.Response, error) {
	if val, ok := c.Cache.Get(url); ok {
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
			c.Cache.Add(url, data, false)
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

func isTokenExpired(tokenString string) (bool, error) {
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
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		// If the server responds with a 401 (unauthorized) code,
		// we need to check whether or not it concerns an expired token.
		tokenExpired, err := isTokenExpired(token)

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
			err = c.NewTokenOrLogout()
			if err != nil {
				// an error return here means we've been logged out
				return nil, err
			}
			token = c.LoggedInUser.Token
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
		c.Cache.DeleteAllStartsWith(url)
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		path, err := getURLResourcePath(url)
		if err != nil {
			slog.Error("could not delete keys with url prefix from cache", slog.String("prefix", url))
		} else {
			// delete all probable cache of this resource
			c.Cache.DeleteAllStartsWith(path)
		}
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

// NewTokenOrLogout attempts to get a new access token from the server
// using the refresh token stored in the current session.
// If successful, the new access token is stored into session.
// If an access token can not be retrieved, the refresh token is
// revoked and the user is logged out of the client.
func (c *Client) NewTokenOrLogout() error {
	accessToken, err := c.GetAccessToken(c.LoggedInUser.RefreshToken)
	if err != nil {
		// If the refresh token is invalid, try to revoke it and,
		// after clearing cache of the current user session, return the error.
		revokeErr := c.RevokeRefreshToken(c.LoggedInUser.RefreshToken)
		c.ClearUserSession()
		revokeResult := "success"
		if revokeErr != nil {
			revokeResult = "failure"
		}
		slog.Info("Failed to get new access token. Attempted to revoke refresh token with result: " + revokeResult)
		return err
	}
	c.LoggedInUser.Token = accessToken
	return nil
}
