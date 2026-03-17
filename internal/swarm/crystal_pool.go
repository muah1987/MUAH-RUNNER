package swarm

import "sync"

// CrystalPool tracks shared token and API-call budgets across swarm agents.
// A sync.RWMutex is used so that concurrent Remaining() reads don't block
// each other; only Allocate() (which mutates state) requires an exclusive lock.
type CrystalPool struct {
mu           sync.RWMutex
TokenBudget  int
TokensUsed   int
APICalls     int
APICallsUsed int
}

func NewCrystalPool(tokenBudget, apiCalls int) *CrystalPool {
return &CrystalPool{
TokenBudget: tokenBudget,
APICalls:    apiCalls,
}
}

func (c *CrystalPool) Allocate(tokens, calls int) bool {
c.mu.Lock()
defer c.mu.Unlock()
if c.TokensUsed+tokens > c.TokenBudget {
return false
}
if c.APICallsUsed+calls > c.APICalls {
return false
}
c.TokensUsed += tokens
c.APICallsUsed += calls
return true
}

// Remaining returns the unused token and API-call budgets.
// It uses a read lock so multiple callers can query concurrently without
// blocking each other or ongoing Allocate calls.
func (c *CrystalPool) Remaining() (int, int) {
c.mu.RLock()
defer c.mu.RUnlock()
return c.TokenBudget - c.TokensUsed, c.APICalls - c.APICallsUsed
}
