// Package client handles pincher-api calls
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	file "github.com/YouWantToPinch/pincher-cli/internal/filemgr"
	"github.com/golang-jwt/jwt/v5"
)

// Client is a struct offering methods used to interact
// with the Pincher REST API. You should not instantiate
// this client directly, and instead use the NewClient()
// function.
type Client struct {
	http.Client
	Cache         Cache
	baseURL       string
	parsedBaseURL *url.URL
	token         string
	RefreshToken  string
}

// NewClient is the proper way to instantiate a client with the sdk.
func NewClient(timeout, cacheInterval time.Duration, baseURL string) (Client, error) {
	c := Client{
		Cache: *NewCache(cacheInterval),
		Client: http.Client{
			Timeout: timeout,
		},
	}
	err := c.SetBaseURL(baseURL)
	if err != nil {
		return Client{}, fmt.Errorf("client.SetBaseURL: %w", err)
	}
	return c, nil
}

// BaseURL returns the BaseURL stored within the client,
// used for API calls.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// APIURL returns the BaseURL stored within the client,
// concatenated with the appropriate api path.
func (c *Client) APIURL() string {
	return c.baseURL + "/api"
}

// SetBaseURL sets the base URL stored within the client,
// used for API calls.
func (c *Client) SetBaseURL(newURL string) error {
	u, err := validateBaseURL(newURL)
	if err != nil {
		return err
	}

	c.baseURL = u.String()
	c.parsedBaseURL = u

	slog.Info("Client Base URL set", slog.String("BaseURL", c.baseURL))
	return nil
}

// ClearCache calls the Clear() function on
// the cache attached to the client and
// forces an early save of the cache file.
func (c *Client) ClearCache() {
	c.Cache.Clear()
	err := c.SaveCacheFile()
	if err != nil {
		slog.Error("could not save cache file: " + err.Error())
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

	c.Cache.Set(loadedCache.Entries)
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
	err = file.WriteAsJSON(c.Cache, cachePath)
	if err != nil {
		return fmt.Errorf(errMsg+"%w", err)
	}
	return nil
}

// Get wrapper for doRequest
func (c *Client) Get(url, token string, out any) (response *http.Response, respFromCache bool, err error) {
	if val, ok := c.Cache.Get(url); ok {
		slog.Info("retrieving requested data from cache", slog.String("URL", url))
		err := json.Unmarshal(val, out)
		if err != nil {
			slog.Error("attempted to pull request from cached data, but failed: " + err.Error())
		} else {
			return nil, true, nil
		}
	}

	resp, err := c.doRequestWithCache(http.MethodGet, url, token, nil, out)
	if err == nil && out != nil {
		data, cacheErr := json.Marshal(out)
		if cacheErr != nil {
			slog.Error(fmt.Sprintf("could not cache response data: %s", cacheErr))
		} else {
			c.Cache.Add(url, data, false)
		}
	}

	return resp, false, err
}

// Post wrapper for doRequestWithCache
func (c *Client) Post(url, token string, payload, out any) (*http.Response, error) {
	return c.doRequestWithCache(http.MethodPost, url, token, payload, out)
}

// Put wrapper for doRequestWithCache
func (c *Client) Put(url, token string, payload any) (*http.Response, error) {
	return c.doRequestWithCache(http.MethodPut, url, token, payload, nil)
}

// Patch wrapper for doRequestWithCache
func (c *Client) Patch(url, token string, payload any) (*http.Response, error) {
	return c.doRequestWithCache(http.MethodPatch, url, token, payload, nil)
}

// Delete wrapper for doRequestWithCache
func (c *Client) Delete(url, token string, payload any) (*http.Response, error) {
	return c.doRequestWithCache(http.MethodDelete, url, token, payload, nil)
}

// Request validates a request before making a call to the API with it,
// adding the client's internal token value to the Authorization header
// as a Bearer token, if it is valid and not empty.
func (c *Client) Request(method, destination string, data, result any) error {
	return c.doRequest(&c.token, method, destination, data, result)
}

// RequestWithToken validates a request before making a call to the API with it,
// but instead of using the client's internal token value as a bearer
// token in the Authorization header, accepts the given token and uses
// that, instead.
//
// This is useful when making calls to the API that may expect a
// Refresh Token rather than an Access Token, so as to get a NEW access
// token to work with.
func (c *Client) RequestWithToken(token *string, method, destination string, data, result any) error {
	return c.doRequest(token, method, destination, data, result)
}

// doRequest validates a request before making a call to the API with it.
func (c *Client) doRequest(token *string, method, destination string, data, result any) error {
	destination, err := c.ResolveURL("/api" + destination)
	if err != nil {
		return err
	}

	reader, contentType, err := c.prepareRequestBody(data)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(context.Background(), method, destination, reader)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", contentType)

	if token != nil && *token != "" {
		request.Header.Set("Authorization", "Bearer "+*token)
	}

	response, err := c.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	err = c.handleResponse(response.StatusCode, response.Body, result)
	if err != nil {
		return fmt.Errorf("for destination %s: %w", destination, err)
	}

	return nil
}

// prepareJSONBody encodes data as JSON
func (c *Client) prepareJSONBody(body any) (io.Reader, string, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, "", fmt.Errorf("json.Marshal: %w", err)
	}

	return bytes.NewReader(data), "application/json", nil
}

func (c *Client) prepareRequestBody(body any) (io.Reader, string, error) {
	if body == nil {
		return http.NoBody, "application/json", nil
	}

	return c.prepareJSONBody(body)
}

// handleResponse processes the API response
func (c *Client) handleResponse(statusCode int, body io.Reader, result any) error {
	switch statusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusOK, http.StatusCreated:
		if result != nil {
			if err := json.NewDecoder(body).Decode(result); err != nil {
				return fmt.Errorf("handleResponse: %w", err)
			}
		}
	default:
		const limit = 1024
		message, _ := io.ReadAll(io.LimitReader(body, limit))
		return fmt.Errorf("bad status code %d: %s", statusCode, message)
	}

	return nil
}

// ResolveURL converts a relative URL to an absolute URL.
// It prefixes relative URLs with the API base URL.
func (c *Client) ResolveURL(destination string) (string, error) {
	destination = strings.TrimSpace(destination)
	if destination == "" {
		return "", fmt.Errorf("destination empty")
	}

	u, err := url.Parse(destination)
	if err != nil {
		return "", fmt.Errorf("parse(destination): %w", err)
	}

	// Reject scheme-less URLs (//host/path) and any provided scheme.
	if u.Scheme != "" || u.Host != "" {
		if sameHostname(u, c.parsedBaseURL) {
			return u.String(), nil
		}
		return "", fmt.Errorf("refusing external URL host %q", u.Host)
	}

	// Path-only (or query/fragment) reference.
	return c.parsedBaseURL.ResolveReference(u).String(), nil
}

func sameHostname(a, b *url.URL) bool {
	// Host may include port; compare case-insensitively.
	return strings.EqualFold(a.Host, b.Host)
}

// doRequestWithCache is a DEPRECATED function that makes an API call
// with a request built from the given arguments, and updates the client's
// internal cache when doing so.
// doRequestWithCache is being phased out in favor of an approach that expects
// package users to deliberately work with the internal cache is needed, rather
// than baking its logic into any logic making API calls.
// doRequestWithCache has been replaced with Client.Request().
func (c *Client) doRequestWithCache(method, url, token string, payload, out any) (*http.Response, error) {
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
		// we need to check whether or not it concerns an expired access token.
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
			err := c.UserTokenRefresh()
			if err != nil {
				// an error return here means the refrsh token was invalid
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

// MakeRequest is a DEPRECATED function that accepts
// arguments for building a request to be used for
// an API call and returns the result.
// In doRequest, it is replaced with prepareRequestBody.
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

// --------------------------------------------
//  HTTP data that can be sent to the REST API
// --------------------------------------------

// AUTH, USERS

type UserCreateData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type (
	UserLoginData  = UserCreateData
	UserUpdateData = UserCreateData
	UserDeleteData = UserCreateData
)

// RESOURCE META VALUES

type MetaData struct {
	Name  string `json:"name"`
	Notes string `json:"notes"`
}

// BUDGETS

type BudgetCreateData struct {
	MetaData
}

type BudgetUpdateData = BudgetCreateData

// ACCOUNTS

const (
	BudgetAccountTypeOnBudget  = "ON_BUDGET"
	BudgetAccountTypeOffBudget = "OFF_BUDGET"
)

type BudgetAccountCreateData struct {
	MetaData
	AccountType string `json:"account_type"`
}

type BudgetAccountUpdateData = BudgetAccountCreateData

type BudgetAccountDeleteData struct {
	DeleteHard bool `json:"delete_hard"`
}

// GROUPS

type BudgetGroupCreateData struct {
	MetaData
}
type BudgetGroupUpdateData = BudgetGroupCreateData

// PAYEES

type BudgetPayeeCreateData struct {
	MetaData
}
type BudgetPayeeUpdateData = BudgetPayeeCreateData

type BudgetPayeeDeleteData struct {
	NewPayeeName string `json:"new_payee_name"`
}

// CATEGORIES

type BudgetCategoryCreateData struct {
	MetaData
	GroupName string `json:"group_name"`
}
type BudgetCategoryUpdateData = BudgetCategoryCreateData

type BudgetCategoryAssignData struct {
	Amount       int64  `json:"amount"`
	ToCategory   string `json:"to_category"`
	FromCategory string `json:"from_category"`
}

// TRANSACTIONS

type BudgetTransactionCreateData struct {
	AccountName         string           `json:"account_name"`
	TransferAccountName string           `json:"transfer_account_name"`
	TransactionDate     string           `json:"transaction_date"`
	PayeeName           string           `json:"payee_name"`
	Notes               string           `json:"notes"`
	Cleared             bool             `json:"is_cleared"`
	Amounts             map[string]int64 `json:"amounts"`
}
type BudgetTransactionUpdateData = BudgetTransactionCreateData
