package client

import (
	"maps"
	"strings"
	"sync"
	"time"
)

type metadata struct {
	DestinationURL string    `json:"destination_url"` // request by which this entry was acquired
	CreatedAt      time.Time `json:"created_at"`      // time at which this entry was created
	Protected      bool      `json:"protected"`       // whether this entry is protected from the reap loop
	Data           []byte    `json:"data"`            // data tied to entry
}

type accountCacheEntry struct {
	account *Account
	metadata
}

type payeeCacheEntry struct {
	payee *Payee
	metadata
}

type groupCacheEntry struct {
	group *Group
	metadata
}

type categoryCacheEntry struct {
	category *Category
	metadata
}

type txnCacheEntry struct {
	transaction *Transaction
	metadata
}

type txnDetailsCacheEntry struct {
	txnDetails *TransactionDetail
	metadata
}

// budgetCache is a subcache that stores its own relevant
// budget resources such as accounts, payees, groups, etc.
type budgetCache struct {
	AccountCache    map[string]*accountCacheEntry    `json:"account_cache"`
	PayeeCache      map[string]*payeeCacheEntry      `json:"payee_cache"`
	GroupCache      map[string]*groupCacheEntry      `json:"group_cache"`
	CategoryCache   map[string]*categoryCacheEntry   `json:"category_cache"`
	TxnCache        map[string]*txnCacheEntry        `json:"transaction_cache"`
	TxnDetailsCache map[string]*txnDetailsCacheEntry `json:"transaction_details_cache"`
	budgetEntry     Budget
	entryMetadata   metadata
}

// GET
func (c *budgetCache) Budget() *Budget {
	return &c.budgetEntry
}

func (c *budgetCache) Account(aID string) *Account {
	return c.AccountCache[aID].account
}

func (c *budgetCache) Payee(pID string) *Payee {
	return c.PayeeCache[pID].payee
}

func (c *budgetCache) Group(gID string) *Group {
	return c.GroupCache[gID].group
}

func (c *budgetCache) Category(cID string) *Category {
	return c.CategoryCache[cID].category
}

func (c *budgetCache) Transaction(tID string) *Transaction {
	return c.TxnCache[tID].transaction
}

func (c *budgetCache) TransactionDetail(tID string) *TransactionDetail {
	return c.TxnDetailsCache[tID].txnDetails
}

type Cache struct {
	mu *sync.Mutex

	// map of budget IDs to a subcache of budget resources
	budgets map[string]*budgetCache

	// whether to update cache entries from relevant API calls
	// related to singleton resources
	trackAPICalls bool

	// whether to update cache entries from relevant API calls
	// related to resource collections
	trackBulkAPICalls bool

	Entries  map[string]metadata `json:"cached_entries"`
	interval time.Duration
}

// ----------------------
//    GETTER FUNCTIONS
// ----------------------

func (c *Cache) Budget(bID string) *Budget {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.budgets[bID]
	if !ok {
		return nil
	}

	return b.Budget()
}

func (c *Cache) Budgets(urlQuery string) []*Budget {
	c.mu.Lock()
	defer c.mu.Unlock()

	var budgets []*Budget
	for _, b := range c.budgets {
		if urlQuery != b.entryMetadata.DestinationURL {
			continue
		}
		if b.entryMetadata.DestinationURL != EndpointBudgets()+urlQuery {
			continue
		}
		budgets = append(budgets, &b.budgetEntry)
	}

	return budgets
}

func (c *Cache) Account(bID, aID string) *Account {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.budgets[bID]
	if !ok {
		return nil
	}

	return b.Account(aID)
}

func (c *Cache) Accounts(bID, urlQuery string) []*Account {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.budgets[bID]
	if !ok {
		return nil
	}

	accounts := make([]*Account, 0, len(b.AccountCache))
	for _, entry := range b.AccountCache {
		if entry.DestinationURL != EndpointBudgetAccounts(bID)+urlQuery {
			continue
		}
		accounts = append(accounts, entry.account)
	}

	return accounts
}

func (c *Cache) Payee(bID, pID string) *Payee {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.budgets[bID]
	if !ok {
		return nil
	}

	return b.Payee(pID)
}

func (c *Cache) Payees(bID, urlQuery string) []*Payee {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.budgets[bID]
	if !ok {
		return nil
	}

	payees := make([]*Payee, 0, len(b.PayeeCache))
	for _, entry := range b.PayeeCache {
		if entry.DestinationURL != EndpointBudgetPayees(bID)+urlQuery {
			continue
		}
		payees = append(payees, entry.payee)
	}

	return payees
}

func (c *Cache) Group(bID, gID string) *Group {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.budgets[bID]
	if !ok {
		return nil
	}

	return b.Group(gID)
}

func (c *Cache) Groups(bID, urlQuery string) []*Group {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.budgets[bID]
	if !ok {
		return nil
	}

	groups := make([]*Group, 0, len(b.GroupCache))
	for _, entry := range b.GroupCache {
		if entry.DestinationURL != EndpointBudgetGroups(bID)+urlQuery {
			continue
		}
		groups = append(groups, entry.group)
	}

	return groups
}

func (c *Cache) Category(bID, cID string) *Category {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.budgets[bID]
	if !ok {
		return nil
	}

	return b.Category(cID)
}

func (c *Cache) Categories(bID, urlQuery string) []*Category {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.budgets[bID]
	if !ok {
		return nil
	}

	categories := make([]*Category, 0, len(b.CategoryCache))
	for _, entry := range b.CategoryCache {
		if entry.DestinationURL != EndpointBudgetCategories(bID)+urlQuery {
			continue
		}
		categories = append(categories, entry.category)
	}

	return categories
}

func (c *Cache) Transaction(bID, tID string) *Transaction {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.budgets[bID]
	if !ok {
		return nil
	}

	return b.Transaction(tID)
}

func (c *Cache) Transactions(bID, urlQuery string) []*Transaction {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.budgets[bID]
	if !ok {
		return nil
	}

	txns := make([]*Transaction, 0, len(b.TxnCache))
	for _, entry := range b.TxnCache {
		if entry.DestinationURL != EndpointBudgetTransactions(bID)+urlQuery {
			continue
		}
		txns = append(txns, entry.transaction)
	}

	return txns
}

func (c *Cache) TransactionDetails(bID, tID string) *TransactionDetail {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.budgets[bID]
	if !ok {
		return nil
	}

	return b.TransactionDetail(tID)
}

func (c *Cache) TransactionsDetails(bID, urlQuery string) []*TransactionDetail {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.budgets[bID]
	if !ok {
		return nil
	}

	txns := make([]*TransactionDetail, 0, len(b.TxnDetailsCache))
	for _, entry := range b.TxnDetailsCache {
		if entry.DestinationURL != EndpointBudgetTransactionsDetails(bID)+urlQuery {
			continue
		}
		txns = append(txns, entry.txnDetails)
	}

	return txns
}

// ----------------------
//    SETTER FUNCTIONS
// ----------------------

func (c *Cache) addBudget(dest, bID string, budget *Budget) {
	if !c.trackAPICalls || budget == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}

	c.budgets[bID] = &budgetCache{
		AccountCache:    map[string]*accountCacheEntry{},
		PayeeCache:      map[string]*payeeCacheEntry{},
		GroupCache:      map[string]*groupCacheEntry{},
		CategoryCache:   map[string]*categoryCacheEntry{},
		TxnCache:        map[string]*txnCacheEntry{},
		TxnDetailsCache: map[string]*txnDetailsCacheEntry{},
		budgetEntry:     *budget,
		entryMetadata: metadata{
			CreatedAt:      time.Now().UTC(),
			DestinationURL: dest,
		},
	}
}

func (c *Cache) addBudgets(dest string, budgets []*Budget) {
	if !c.trackBulkAPICalls || budgets == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, b := range budgets {
		c.budgets[b.ID.String()] = &budgetCache{
			AccountCache:    map[string]*accountCacheEntry{},
			PayeeCache:      map[string]*payeeCacheEntry{},
			GroupCache:      map[string]*groupCacheEntry{},
			CategoryCache:   map[string]*categoryCacheEntry{},
			TxnCache:        map[string]*txnCacheEntry{},
			TxnDetailsCache: map[string]*txnDetailsCacheEntry{},
			budgetEntry:     *b,
			entryMetadata: metadata{
				CreatedAt:      time.Now().UTC(),
				DestinationURL: dest,
			},
		}
	}
}

func (c *Cache) deleteBudget(bID string) {
	if !c.trackAPICalls {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}

	delete(c.budgets, bID)
}

func (c *Cache) addAccount(dest, bID string, account *Account) {
	if !c.trackAPICalls || account == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}

	c.budgets[bID].AccountCache[account.ID.String()] = &accountCacheEntry{
		account: account,
		metadata: metadata{
			CreatedAt:      time.Now().UTC(),
			DestinationURL: dest,
		},
	}
}

func (c *Cache) addAccounts(dest, bID string, accounts []*Account) {
	if !c.trackBulkAPICalls || accounts == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}
	for _, a := range accounts {
		if a == nil {
			continue
		}
		c.budgets[bID].AccountCache[a.ID.String()] = &accountCacheEntry{
			account: a,
			metadata: metadata{
				CreatedAt:      time.Now().UTC(),
				DestinationURL: dest,
			},
		}
	}
}

func (c *Cache) deleteAccount(bID, aID string) {
	if !c.trackAPICalls {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}

	delete(c.budgets[bID].AccountCache, aID)
}

func (c *Cache) addPayee(dest, bID string, payee *Payee) {
	if !c.trackAPICalls || payee == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}

	c.budgets[bID].PayeeCache[payee.ID.String()] = &payeeCacheEntry{
		payee: payee,
		metadata: metadata{
			CreatedAt:      time.Now().UTC(),
			DestinationURL: dest,
		},
	}
}

func (c *Cache) addPayees(dest, bID string, payees []*Payee) {
	if !c.trackBulkAPICalls || payees == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}
	for _, p := range payees {
		if p == nil {
			continue
		}
		c.budgets[bID].PayeeCache[p.ID.String()] = &payeeCacheEntry{
			payee: p,
			metadata: metadata{
				CreatedAt:      time.Now().UTC(),
				DestinationURL: dest,
			},
		}
	}
}

func (c *Cache) deletePayee(bID, pID string) {
	if !c.trackAPICalls {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}

	delete(c.budgets[bID].PayeeCache, pID)
}

func (c *Cache) addGroup(dest, bID string, group *Group) {
	if !c.trackAPICalls || group == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}

	c.budgets[bID].GroupCache[group.ID.String()] = &groupCacheEntry{
		group: group,
		metadata: metadata{
			CreatedAt:      time.Now().UTC(),
			DestinationURL: dest,
		},
	}
}

func (c *Cache) addGroups(dest, bID string, groups []*Group) {
	if !c.trackBulkAPICalls || groups == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}
	for _, g := range groups {
		if g == nil {
			continue
		}
		c.budgets[bID].GroupCache[g.ID.String()] = &groupCacheEntry{
			group: g,
			metadata: metadata{
				CreatedAt:      time.Now().UTC(),
				DestinationURL: dest,
			},
		}
	}
}

func (c *Cache) deleteGroup(bID, gID string) {
	if !c.trackAPICalls {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}

	delete(c.budgets[bID].GroupCache, gID)
}

func (c *Cache) addCategory(dest, bID string, category *Category) {
	if !c.trackAPICalls || category == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}

	c.budgets[bID].CategoryCache[category.ID.String()] = &categoryCacheEntry{
		category: category,
		metadata: metadata{
			CreatedAt:      time.Now().UTC(),
			DestinationURL: dest,
		},
	}
}

func (c *Cache) addCategories(dest, bID string, categories []*Category) {
	if !c.trackBulkAPICalls || categories == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}
	for _, cat := range categories {
		if cat == nil {
			continue
		}
		c.budgets[bID].CategoryCache[cat.ID.String()] = &categoryCacheEntry{
			category: cat,
			metadata: metadata{
				CreatedAt:      time.Now().UTC(),
				DestinationURL: dest,
			},
		}
	}
}

func (c *Cache) deleteCategory(bID, cID string) {
	if !c.trackAPICalls {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}

	delete(c.budgets[bID].CategoryCache, cID)
}

func (c *Cache) addTxn(dest, bID string, txn *Transaction) {
	if !c.trackAPICalls || txn == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}

	c.budgets[bID].TxnCache[txn.ID.String()] = &txnCacheEntry{
		transaction: txn,
		metadata: metadata{
			CreatedAt:      time.Now().UTC(),
			DestinationURL: dest,
		},
	}
}

func (c *Cache) addTxns(dest, bID string, txns []*Transaction) {
	if !c.trackBulkAPICalls || txns == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}
	for _, t := range txns {
		if t == nil {
			continue
		}
		c.budgets[bID].TxnCache[t.ID.String()] = &txnCacheEntry{
			transaction: t,
			metadata: metadata{
				CreatedAt:      time.Now().UTC(),
				DestinationURL: dest,
			},
		}
	}
}

func (c *Cache) deleteTxn(bID, tID string) {
	if !c.trackAPICalls {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}

	delete(c.budgets[bID].TxnCache, tID)
}

func (c *Cache) addTxnDetails(dest, bID string, txn *TransactionDetail) {
	if !c.trackAPICalls || txn == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}

	c.budgets[bID].TxnDetailsCache[txn.ID.String()] = &txnDetailsCacheEntry{
		txnDetails: txn,
		metadata: metadata{
			CreatedAt:      time.Now().UTC(),
			DestinationURL: dest,
		},
	}
}

func (c *Cache) addTxnsDetails(dest, bID string, txns []*TransactionDetail) {
	if !c.trackBulkAPICalls || txns == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}
	for _, t := range txns {
		if t == nil {
			continue
		}
		c.budgets[bID].TxnDetailsCache[t.ID.String()] = &txnDetailsCacheEntry{
			txnDetails: t,
			metadata: metadata{
				CreatedAt:      time.Now().UTC(),
				DestinationURL: dest,
			},
		}
	}
}

func (c *Cache) deleteTxnsDetails(bID, tID string) {
	if !c.trackAPICalls {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.budgets[bID]; !ok {
		return
	}

	delete(c.budgets[bID].TxnDetailsCache, tID)
}

// ---------------------------
//  DEPRECATED CACHE FUNCTIONS
// ---------------------------

func (c *Cache) Set(entries map[string]metadata) {
	c.mu.Lock()
	defer c.mu.Unlock()

	maps.Copy(c.Entries, entries)
}

func (c *Cache) Add(key string, value []byte, protect bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Entries[key] = metadata{
		CreatedAt: time.Now().UTC(),
		Data:      value,
		Protected: protect,
	}
}

func (c *Cache) Get(key string) (entryData []byte, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.Entries[key]
	if !ok {
		return nil, false
	}

	return val.Data, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.reap()
	}
}

func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO: I know there's a way one can simplify this with
	// an interface instead of all these multiple loops
	// FOR EACH subcache, but I was getting errors when
	// trying...

	for _, bCache := range c.budgets {
		for aID, entry := range bCache.AccountCache {
			if !entry.Protected && (time.Since(entry.CreatedAt) > c.interval) {
				delete(bCache.AccountCache, aID)
			}
		}
		for pID, entry := range bCache.PayeeCache {
			if !entry.Protected && (time.Since(entry.CreatedAt) > c.interval) {
				delete(bCache.PayeeCache, pID)
			}
		}
		for gID, entry := range bCache.GroupCache {
			if !entry.Protected && (time.Since(entry.CreatedAt) > c.interval) {
				delete(bCache.GroupCache, gID)
			}
		}
		for cID, entry := range bCache.CategoryCache {
			if !entry.Protected && (time.Since(entry.CreatedAt) > c.interval) {
				delete(bCache.CategoryCache, cID)
			}
		}
		for tID, entry := range bCache.TxnCache {
			if !entry.Protected && (time.Since(entry.CreatedAt) > c.interval) {
				delete(bCache.TxnCache, tID)
			}
		}
		for tID, entry := range bCache.TxnDetailsCache {
			if !entry.Protected && (time.Since(entry.CreatedAt) > c.interval) {
				delete(bCache.TxnDetailsCache, tID)
			}
		}
	}

	for key, val := range c.Entries {
		if !val.Protected && (time.Since(val.CreatedAt) > c.interval) {
			delete(c.Entries, key)
		}
	}
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.Entries, key)
}

// DeleteAllStartsWith removes all cached entries whose
// key starts with the given prefix.
// This can be used to invalidate all entries related
// to a resource that may have had an instance updated or removed.
func (c *Cache) DeleteAllStartsWith(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.Entries {
		if strings.HasPrefix(key, prefix) {
			delete(c.Entries, key)
		}
	}
}

// Clear deletes all cached entries.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	clear(c.Entries)
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		mu:                &sync.Mutex{},
		interval:          interval,
		trackAPICalls:     true,
		trackBulkAPICalls: true,
		Entries:           make(map[string]metadata),
		budgets:           map[string]*budgetCache{},
	}

	go cache.reapLoop()

	return &cache
}
