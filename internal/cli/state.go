package cli

import (
	"fmt"
	"log/slog"

	"github.com/YouWantToPinch/pincher-cli/internal/config"
	file "github.com/YouWantToPinch/pincher-cli/internal/filemgr"
	pgo "github.com/YouWantToPinch/pincher-sdk-go/pinchergo"
)

// State represents the full state of the CLI during a session,
// regardless of user login. It knows the full context of all
// systems connected.
type State struct {
	CmdQueue chan<- string
	DoneChan *chan bool
	Logger   *Logger
	Config   *config.Config
	Client   *pgo.Client
	Session  *cliSession
	styles   *styles
}

// GetBudget goes through the Client to retrieve a
// budget belonging to the active user. It first
// attempts to pull from cache, then making an API
// call if it is unable to do so.
func (s *State) GetBudget(bID string) (budget *pgo.Budget, err error) {
	budget = s.Client.Cache.Budget(bID)
	if budget != nil {
		return budget, nil
	}
	budget, err = s.Client.Budget(bID)
	return budget, err
}

// GetBudgets goes through the Client to retrieve a list of
// budgets belonging to the active user. It first attempts
// to pull from cache, then making an API call if it is
// unable to do so.
func (s *State) GetBudgets(bID, urlQuery string) (budgets []*pgo.Budget, err error) {
	budgets = s.Client.Cache.Budgets(urlQuery)
	if budgets != nil {
		return budgets, nil
	}
	budgets, err = s.Client.Budgets(bID, urlQuery)
	return budgets, err
}

// GetAccount goes through the Client to retrieve an
// account by ID belonging to the budget which corresponds
// to the given budget ID. It first attempts to pull from
// cache, then making an API call if it is unable to do so.
func (s *State) GetAccount(bID, aID string) (account *pgo.Account, err error) {
	account = s.Client.Cache.Account(bID, aID)
	if account != nil {
		return account, nil
	}
	account, err = s.Client.BudgetAccount(bID, aID)
	return account, err
}

// GetAccounts goes through the Client to retrieve a list of
// accounts belonging to budget which corresponds to the given
// budget ID. It first attempts to pull from cache, then
// making an API call if it is unable to do so.
func (s *State) GetAccounts(bID, urlQuery string) (accounts []*pgo.Account, err error) {
	accounts = s.Client.Cache.Accounts(bID, urlQuery)
	if len(accounts) > 0 {
		return accounts, nil
	}
	accounts, err = s.Client.BudgetAccounts(bID, urlQuery)
	return accounts, err
}

// GetPayee goes through the Client to retrieve an
// payee by ID belonging to the budget which corresponds
// to the given budget ID. It first attempts to pull from
// cache, then making an API call if it is unable to do so.
func (s *State) GetPayee(bID, pID string) (payee *pgo.Payee, err error) {
	payee = s.Client.Cache.Payee(bID, pID)
	if payee != nil {
		return payee, nil
	}
	payee, err = s.Client.BudgetPayee(bID, pID)
	return payee, err
}

// GetPayees goes through the Client to retrieve a list of
// payees belonging to budget which corresponds to the given
// budget ID. It first attempts to pull from cache, then
// making an API call if it is unable to do so.
func (s *State) GetPayees(bID, urlQuery string) (payees []*pgo.Payee, err error) {
	payees = s.Client.Cache.Payees(bID, urlQuery)
	if len(payees) > 0 {
		return payees, nil
	}
	payees, err = s.Client.BudgetPayees(bID, urlQuery)
	return payees, err
}

// GetGroup goes through the Client to retrieve an
// group by ID belonging to the budget which corresponds
// to the given budget ID. It first attempts to pull from
// cache, then making an API call if it is unable to do so.
func (s *State) GetGroup(bID, gID string) (group *pgo.Group, err error) {
	group = s.Client.Cache.Group(bID, gID)
	if group != nil {
		return group, nil
	}
	group, err = s.Client.BudgetGroup(bID, gID)
	return group, err
}

// GetGroups goes through the Client to retrieve a list of
// groups belonging to budget which corresponds to the given
// budget ID. It first attempts to pull from cache, then
// making an API call if it is unable to do so.
func (s *State) GetGroups(bID, urlQuery string) (groups []*pgo.Group, err error) {
	groups = s.Client.Cache.Groups(bID, urlQuery)
	if len(groups) > 0 {
		return groups, nil
	}
	groups, err = s.Client.BudgetGroups(bID, urlQuery)
	return groups, err
}

// GetCategory goes through the Client to retrieve an
// category by ID belonging to the budget which corresponds
// to the given budget ID. It first attempts to pull from
// cache, then making an API call if it is unable to do so.
func (s *State) GetCategory(bID, cID string) (category *pgo.Category, err error) {
	category = s.Client.Cache.Category(bID, cID)
	if category != nil {
		return category, nil
	}
	category, err = s.Client.BudgetCategory(bID, cID)
	return category, err
}

// GetCategories goes through the Client to retrieve a list of
// categories belonging to budget which corresponds to the given
// budget ID. It first attempts to pull from cache, then
// making an API call if it is unable to do so.
func (s *State) GetCategories(bID, urlQuery string) (categories []*pgo.Category, err error) {
	categories = s.Client.Cache.Categories(bID, urlQuery)
	if len(categories) > 0 {
		return categories, nil
	}
	categories, err = s.Client.BudgetCategories(bID, urlQuery)
	return categories, err
}

// GetTxn goes through the Client to retrieve an
// transaction by ID belonging to the budget which
// corresponds to the given budget ID. It first
// attempts to pull from cache, then making an API
// call if it is unable to do so.
func (s *State) GetTxn(bID, tID string) (txn *pgo.Transaction, err error) {
	txn = s.Client.Cache.Transaction(bID, tID)
	if txn != nil {
		return txn, nil
	}
	txn, err = s.Client.BudgetTransaction(bID, tID)
	return txn, err
}

// GetTxns goes through the Client to retrieve a list of
// transactions belonging to budget which corresponds to the given
// budget ID. It first attempts to pull from cache, then
// making an API call if it is unable to do so.
func (s *State) GetTxns(bID, urlQuery string) (txns []*pgo.Transaction, err error) {
	txns = s.Client.Cache.Transactions(bID, urlQuery)
	if len(txns) > 0 {
		return txns, nil
	}
	txns, err = s.Client.BudgetTransactions(bID, urlQuery)
	return txns, err
}

// GetTxnDetails goes through the Client to retrieve an
// transaction's details by ID belonging to the budget
// which corresponds to the given budget ID. It first
// attempts to pull from cache, then making an API call
// if it is unable to do so.
func (s *State) GetTxnDetails(bID, tID string) (txn *pgo.TransactionDetail, err error) {
	txn = s.Client.Cache.TransactionDetails(bID, tID)
	if txn != nil {
		return txn, nil
	}
	txn, err = s.Client.BudgetTransactionDetails(bID, tID)
	return txn, err
}

// GetTxnsDetails goes through the Client to retrieve a list of
// transaction details belonging to budget which corresponds to
// the given budget ID. It first attempts to pull from cache, then
// making an API call if it is unable to do so.
func (s *State) GetTxnsDetails(bID, urlQuery string) (txns []*pgo.TransactionDetail, err error) {
	txns = s.Client.Cache.TransactionsDetails(bID, urlQuery)
	if len(txns) > 0 {
		return txns, nil
	}
	txns, err = s.Client.BudgetTransactionsDetails(bID, urlQuery)
	return txns, err
}

// ClearCache calls the Clear() function on
// the cache attached to the pgo.and
// forces an early save of the cache file.
func (s *State) ClearCache() {
	s.Client.Cache.Clear()
	err := s.SaveCacheFile()
	if err != nil {
		slog.Error("could not save cache file: " + err.Error())
	}
}

// LoadCacheFile looks for a file with the given name within
// the user cache directory and attempts to load it into the cache.
func (s *State) LoadCacheFile() error {
	const errMsg = "could not load cache: "
	cachePath, err := file.GetCacheFilepath("cache.json")
	if err != nil {
		return fmt.Errorf(errMsg+"%w", err)
	}

	loadedCache, err := file.ReadJSONFromFile[pgo.Cache](cachePath)
	if err != nil {
		return fmt.Errorf(errMsg+"%w", err)
	}

	s.Client.Cache.Set(loadedCache.Entries)
	return nil
}

// SaveCacheFile saves the current cache to a local file
// with the given name under the user cache directory.
func (s *State) SaveCacheFile() error {
	const errMsg = "could not save cache: "
	cachePath, err := file.GetCacheFilepath("cache.json")
	if err != nil {
		return fmt.Errorf(errMsg+"%w", err)
	}
	err = file.WriteAsJSON(s.Client.Cache, cachePath)
	if err != nil {
		return fmt.Errorf(errMsg+"%w", err)
	}
	return nil
}

func (s *State) NewSession() {
	s.Session = &cliSession{}
	s.Session.Init()
}

// GetPrompt returns the proper input
// prompt to print to the command line,
// dependent on user activity during
// this session.
func (s *State) GetPrompt() string {
	if s.styles != nil {
		return s.getStyledPrompt()
	}
	return s.getUnstyledPrompt()
}

func (s *State) getDiv(styled bool) string {
	div := "___________________________"
	if styled && s.styles != nil {
		return s.styles.Green.Render(div)
	}
	return div
}

// getUnstyledPrompt returns the pincher-cli prompt without pretty colors :(
func (s *State) getUnstyledPrompt() string {
	if s.Session.ActiveUser.Username == "" {
		return "pin¢her > "
	} else if s.Session.ActiveBudget.Name == "" {
		return fmt.Sprintf("p¢/%s > ", s.Session.ActiveUser.Username)
	} else {
		return fmt.Sprintf("p¢/%s[%s] > ", s.Session.ActiveUser.Username, s.Session.ActiveBudget.Name)
	}
}

// getStyledPrompt returns the pincher-cli prompt with pretty colors :)
func (s *State) getStyledPrompt() string {
	cent := s.styles.Orange.Render("¢")
	chev := s.styles.White.Render(" > ")

	if s.Session.ActiveUser.Username == "" {
		pin := s.styles.White.Render("pin")
		her := s.styles.White.Render("her")
		return pin + cent + her + chev
	} else {
		p := s.styles.White.Render("p")
		slash := s.styles.White.Render("/")

		if s.Session.ActiveBudget.Name == "" {
			return p + cent + slash + s.styles.White.Render(s.Session.ActiveUser.Username) + chev
		} else {
			lbr := s.styles.Green.Render("[")
			rbr := s.styles.Green.Render("]")
			return p + cent + slash + s.styles.White.Render(s.Session.ActiveUser.Username) + lbr + s.styles.White.Render(s.Session.ActiveBudget.Name) + rbr + chev
		}
	}
}

// cliSession represents the state of the CLI in regard to
// a logged-in user.
type cliSession struct {
	ActiveUser      pgo.User
	ActiveBudget    pgo.Budget
	CommandRegistry *commandRegistry
}

// Init preregisters all commands to the internal command registry.
func (s *cliSession) Init() {
	// initialize maps under command registry
	s.CommandRegistry = &commandRegistry{
		handlers: make(map[string]*cmdHandler),
		registry: make(map[string]registrationStatus),
	}
	// preregister ALL commands
	s.CommandRegistry.batchRegistration(makeBaseCommandHandlers(), Preregistered)
	s.CommandRegistry.preregister(makeBudgetCommandHandler())
	s.CommandRegistry.batchRegistration(makeResourceCommandHandlers(), Preregistered)
	// register base commands
	s.CommandRegistry.batchRegistration(makeBaseCommandHandlers(), Registered)
}

func (s *cliSession) OnLogin(user pgo.User) {
	// register commands that require login
	s.CommandRegistry.register("budget")
	s.ActiveUser = user
	fmt.Printf("Logged in as user: %s\n", s.ActiveUser.Username)
}

func (s *cliSession) OnViewBudget() {
	// register commands that require viewing a budget
	s.CommandRegistry.batchRegistration(makeResourceCommandHandlers(), Registered)
}

func (s *cliSession) OnLogout() {
	// deregister commands
	s.ActiveUser = pgo.User{}
	s.CommandRegistry.deregisterNonBaseCommands()
}
