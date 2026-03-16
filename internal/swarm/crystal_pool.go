package swarm

import "sync"

type CrystalPool struct {
mu             sync.Mutex
TokenBudget    int
TokensUsed     int
APICalls       int
APICallsUsed   int
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

func (c *CrystalPool) Remaining() (int, int) {
c.mu.Lock()
defer c.mu.Unlock()
return c.TokenBudget - c.TokensUsed, c.APICalls - c.APICallsUsed
}
