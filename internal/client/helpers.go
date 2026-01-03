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
			c.Cache.Add(url, data)
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
	var req *http.Request
	var resp *http.Response
	var err error
	retry := true
	for {
		// try to make the new request.
		req, err = c.MakeRequest(method, url, token, payload)
		if err != nil {
			return nil, err
		}

		// try to send the request.
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return resp, err
		}
		if resp.StatusCode != http.StatusUnauthorized {
			break
		} else if retry {
			tokenExpired, err := isTokenExpired(token)
			if token == "" || (!tokenExpired && err == nil) || (strings.Contains(url, "/refresh") || strings.Contains(url, "revoke")) {
				// If the token is empty or otherwise note needed, that isn't an issue;
				// If the access token is not expired, that isn't an issue;
				// We also need to avoid infinite calls to get/revoke refresh tokens;
				// Break out with the 401.
				break
			}
			// Try to get a new access token if it is invalid or we got an error.
			accessToken, err := c.GetAccessToken(c.LoggedInUser.RefreshToken)
			if err != nil {
				// If the refresh token is invalid, try to revoke it and,
				// after clearing cache of the current user session, return any error.
				err = c.RevokeRefreshToken(c.LoggedInUser.RefreshToken)
				c.LogoutUser()
				return nil, err
			}
			c.LoggedInUser.Token = accessToken
			token = c.LoggedInUser.Token
			retry = false
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
		c.Cache.Delete(url)
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		path, err := getURLResourcePath(url)
		if err != nil {
			slog.Error("could not delete key from url entry cache", slog.String("key", url))
		}
		c.Cache.Delete(path)
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
