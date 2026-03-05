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
	autoRefresh   bool // whether or not to retrieve a new token on expiration, for each API call
}

// NewClient is the proper way to instantiate a client with the sdk.
func NewClient(baseURL string, timeout, cacheInterval time.Duration, autoRefresh bool) (Client, error) {
	c := Client{
		Cache: *NewCache(cacheInterval),
		Client: http.Client{
			Timeout: timeout,
		},
		autoRefresh: autoRefresh,
	}
	err := c.SetBaseURL(baseURL)
	if err != nil {
		return Client{}, fmt.Errorf("client.SetBaseURL: %w", err)
	}
	return c, nil
}

// NewClientWithDefaults returns an auto-refreshing client
// with a 10-second timeout and 5-minute cacheInterval.
//
// Its baseURL is set to the value of defaultBaseURL for
// hosting locally.
func NewClientWithDefaults() (Client, error) {
	c, err := NewClient(
		defaultBaseURL,
		time.Second*10,
		time.Minute*5,
		true,
	)
	if err != nil {
		return Client{}, err
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

	var response *http.Response
	destUsesAccessToken := !strings.Contains(destination, "/refresh") && !strings.Contains(destination, "/revoke")
	// avoid an infinite loop by filtering out endpoints that expect
	// a refresh token in place of an access token
	shouldRetry := destUsesAccessToken
	for {
		// try to send the request.
		response, err = c.Do(request)
		if err != nil {
			return nil
		}

		// If the client is not directed to automatically refres haccess tokens,
		// then we won't try in the first place.
		if !c.autoRefresh {
			break
		}

		// If the server responds with a 401 (unauthorized) code,
		// we need to check whether or not it concerns an expired access token.
		tokenExpired, err := checkTokenExpired(*token)

		// if it isn't a 401, we don't care to retry with a new access token
		if response.StatusCode != http.StatusUnauthorized ||
			// if no token was required for the request, we don't care
			*token == "" ||
			// if the token is not expired, the 401 has nothing to do with tokens
			(!tokenExpired && err == nil) {
			break
		} else if shouldRetry {

			// Try to get a new access token if it is invalid or we got an error.
			// If that's not possible, log out the user, as their session must therefore be invalid.
			err := c.UserTokenRefresh()
			if err != nil {
				// an error return here means the refrsh token was invalid
				return err
			}
			// update this loop with the new access token
			token = &c.token
			shouldRetry = false
		} else {
			break
		}
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

// --------------
//  Helpers
// --------------

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
