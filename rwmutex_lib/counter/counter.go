package counter

import "sync"

// 一个计数器
type Counter struct {
	mu    sync.RWMutex
	count uint64
}


// Incr 写保护
func (c *Counter) Incr() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

// 读保护
func (c *Counter) Count() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.count
}
