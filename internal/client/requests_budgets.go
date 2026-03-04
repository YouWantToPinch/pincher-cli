package client

import (
	"net/http"
)

func (c *Client) BudgetCreate(bID string, data BudgetCreateData) error {
	endpoint := EndpointBudgets()
	var budget *Budget
	err := c.Request(http.MethodPost, endpoint, data, &budget)
	c.Cache.addBudget(bID, budget)
	return err
}

type budgetContainer struct {
	Budgets []*Budget `json:"data"`
}

func (c *Client) Budget(bID string) (budget *Budget, err error) {
	endpoint := EndpointBudget(bID)
	err = c.Request(http.MethodGet, endpoint, nil, &budget)
	c.Cache.addBudget(bID, budget)
	return budget, err
}

func (c *Client) Budgets(bID, urlQuery string) (budgets []*Budget, err error) {
	endpoint := EndpointBudgets() + urlQuery
	var container budgetContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	c.Cache.addBudgets(container.Budgets)
	return container.Budgets, err
}

func (c *Client) BudgetReport(bID, mID string) (report *MonthReport, err error) {
	endpoint := EndpointBudgetMonth(bID, mID)
	err = c.Request(http.MethodGet, endpoint, nil, &report)
	return report, err
}

func (c *Client) BudgetUpdate(bID string, data BudgetUpdateData) error {
	endpoint := EndpointBudget(bID)
	err := c.Request(http.MethodPut, endpoint, data, nil)
	if err != nil {
		_, _ = c.Budget(bID)
	}

	return err
}

func (c *Client) BudgetDelete(bID string) error {
	endpoint := EndpointBudget(bID)
	err := c.Request(http.MethodDelete, endpoint, nil, nil)
	return err
}
