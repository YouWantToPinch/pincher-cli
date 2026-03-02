package client

import (
	"fmt"
	"log"
)

/*
	This file contains all the endpoints used in this library for the Pincher API.
	The naming scheme of the constants and methods has rules:

	Constants:
	 - Prefix with "URL"
	 - Follow a hierarchical structure that build on top of each other
	 - Use plural/singular forms to somewhat reflect the relationship between resources;
		- For example, "EndpointChannelAckMessage" relates to a single Channel and a single Message object

	Functions:
	 - Prefix with "Endpoint"
	 - Used to generate a URL for a specific resource
	 - Follow the same hierarchical structure as the constants
*/

/* These base URLs are used by the Client.Request method */
var (
	baseURL = "http://localhost:8080"

	parsedBase = mustParseURL(baseURL + "/api")
)

func BaseURL() string {
	return baseURL
}

func APIURL() string {
	return baseURL + "/api"
}

// SetBaseURL sets the base URL for the API.
func SetBaseURL(newURL string) error {
	u, err := validateBaseURL(newURL)
	if err != nil {
		return err
	}

	baseURL = u.String()
	parsedBase = u

	log.Printf("Base URL set to %s", baseURL)
	return nil
}

const sVerb = "/%s"

const (
	URLHealthz = "/healthz"

	URLAdmin           = "/admin"
	URLAdminUsers      = URLAdmin + "/users"
	URLAdminUsersCount = URLAdminUsers + "/count"

	URLUserTokenLogin   = "/login"
	URLUserTokenRefresh = "/refresh"
	URLUserTokenRevoke  = "/revoke"

	URLUsers = "/users"

	URLBudgets                   = "/budgets"
	URLBudget                    = URLBudgets + sVerb
	URLBudgetCapital             = URLBudget + "/capital"
	URLBudgetMembers             = URLBudget + "/members"
	URLBudgetMember              = URLBudgetMembers + sVerb
	URLBudgetGroups              = URLBudget + "/groups"
	URLBudgetGroup               = URLBudgetGroups + sVerb
	URLBudgetCategories          = URLBudget + "/categories"
	URLBudgetCategory            = URLBudgetCategories + sVerb
	URLBudgetPayees              = URLBudget + "/payees"
	URLBudgetPayee               = URLBudgetPayees + sVerb
	URLBudgetAccounts            = URLBudget + "/accounts"
	URLBudgetAccount             = URLBudgetAccounts + sVerb
	URLBudgetAccountCapital      = URLBudgetAccount + "/capital"
	URLBudgetTransactions        = URLBudget + "/transactions"
	URLBudgetTransactionsDetails = URLBudgetTransactions + "/details"
	URLBudgetTransaction         = URLBudgetTransactions + sVerb
	URLBudgetTransactionDetails  = URLBudgetTransaction + "/details"
	URLBudgetTransactionSplits   = URLBudgetTransaction + "/splits"
	URLBudgetMonths              = URLBudget + "/months"
	URLBudgetMonth               = URLBudgetMonths + sVerb
	URLBudgetMonthCategories     = URLBudgetMonth + "/categories"
	URLBudgetMonthCategory       = URLBudgetMonthCategories + sVerb
)

func EndpointLogin() string {
	return URLUserTokenLogin
}

func EndpointRefresh() string {
	return URLUserTokenRefresh
}

func EndpointRevoke() string {
	return URLUserTokenRevoke
}

func EndpointUsers() string {
	return URLUsers
}

func EndpointBudgets() string {
	return URLBudgets
}

func EndpointBudget(bID string) string {
	return fmt.Sprintf(URLBudget, bID)
}

func EndpointBudgetCapital(bID string) string {
	return fmt.Sprintf(URLBudgetCapital, bID)
}

func EndpointBudgetMembers(bID string) string {
	return fmt.Sprintf(URLBudgetMembers, bID)
}

func EndpointBudgetMember(bID, mID string) string {
	return fmt.Sprintf(URLBudgetMember, bID, mID)
}

func EndpointBudgetGroups(bID string) string {
	return fmt.Sprintf(URLBudgetGroups, bID)
}

func EndpointBudgetGroup(bID, gID string) string {
	return fmt.Sprintf(URLBudgetGroup, bID, gID)
}

func EndpointBudgetCategories(bID string) string {
	return fmt.Sprintf(URLBudgetCategories, bID)
}

func EndpointBudgetCategory(bID, cID string) string {
	return fmt.Sprintf(URLBudgetCategory, bID, cID)
}

func EndpointBudgetPayees(bID string) string {
	return fmt.Sprintf(URLBudgetPayees, bID)
}

func EndpointBudgetPayee(bID, pID string) string {
	return fmt.Sprintf(URLBudgetPayee, bID, pID)
}

func EndpointBudgetAccounts(bID string) string {
	return fmt.Sprintf(URLBudgetAccounts, bID)
}

func EndpointBudgetAccount(bID, aID string) string {
	return fmt.Sprintf(URLBudgetAccount, bID, aID)
}

func EndpointBudgetAccountCapital(bID, aID string) string {
	return fmt.Sprintf(URLBudgetAccountCapital, bID, aID)
}

func EndpointBudgetTransactions(bID string) string {
	return fmt.Sprintf(URLBudgetTransactions, bID)
}

func EndpointBudgetTransaction(bID, tID string) string {
	return fmt.Sprintf(URLBudgetTransaction, bID, tID)
}

func EndpointBudgetTransactionsDetails(bID string) string {
	return fmt.Sprintf(URLBudgetTransactionsDetails, bID)
}

func EndpointBudgetTransactionDetails(bID, tID string) string {
	return fmt.Sprintf(URLBudgetTransactionDetails, bID, tID)
}

func EndpointBudgetTransactionSplits(bID, tID string) string {
	return fmt.Sprintf(URLBudgetTransactionDetails, bID, tID)
}

func EndpointBudgetMonth(bID, mID string) string {
	return fmt.Sprintf(URLBudgetMonth, bID, mID)
}

func EndpointBudgetMonthCategories(bID, mID string) string {
	return fmt.Sprintf(URLBudgetMonthCategories, bID, mID)
}

func EndpointBudgetMonthCategory(bID, mID, cID string) string {
	return fmt.Sprintf(URLBudgetMonthCategory, bID, mID, cID)
}
