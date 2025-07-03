// This package implements a simple threadsafe generic TTL Cache, without LRU mechanism
// deletions are made via channel, therefor the Close() methode.
// As the deletion is async, it's hardly syncronised to all read/write operations and
// secured via an extra isEvicted-check
package ttlcache

import (
	"sync"
	"time"
)

// Cache the simple cache
type Cache[K comparable, V any] struct {
	lock           sync.RWMutex
	items          map[K]entry[V]
	ttl            time.Duration
	deletions      chan K
	autodelete     *time.Ticker
	doneAutodelete chan bool
}

// entry one entry in the map
type entry[V any] struct {
	value     V
	expiresAt time.Time
}

// Option for the functional options (for further development)
type Option[K comparable, V any] func(c *Cache[K, V])

// New creating a new TTL Cache
func New[K comparable, V any](opts ...Option[K, V]) *Cache[K, V] {
	c := &Cache[K, V]{
		lock:      sync.RWMutex{},
		items:     make(map[K]entry[V]),
		ttl:       0,
		deletions: make(chan K, 100),
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.ttl > 0 {
		go func() {
			for k := range c.deletions {
				c.deleteEvicted(k)
			}
		}()
	}
	return c
}

// WithTTL setting the time to life default for all entries
// TTL with 0 will lead into no default TTL, but you can set a
// specifig TTL for every entry
func WithTTL[K comparable, V any](d time.Duration) Option[K, V] {
	return func(c *Cache[K, V]) {
		c.ttl = d
	}
}

// WithNoTTL means there is never set an TTL to any entries.
// mainly this will disable the deletion channel and go func
func WithNoTTL[K comparable, V any]() Option[K, V] {
	return func(c *Cache[K, V]) {
		c.ttl = -1
	}
}

// WithAutoDeletion will start a timer and automatically delete evicted entries. The timer will run every d duration
func WithAutoDeletion[K comparable, V any](d time.Duration) Option[K, V] {
	return func(c *Cache[K, V]) {
		c.autodelete = time.NewTicker(d)
		c.doneAutodelete = make(chan bool)

		go func() {
			for {
				select {
				case <-c.doneAutodelete:
					return
				case <-c.autodelete.C:
					c.DeleteEvicted()
				}
			}
		}()
	}
}

// Stop will stop the auto delete ticker, if present
func (c *Cache[K, V]) Stop() {
	if c.autodelete != nil {
		c.autodelete.Stop()
		c.doneAutodelete <- true
	}
}

// Get getting a value
func (c *Cache[K, V]) Get(key K) (*V, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	e, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if c.isEvicted(e) {
		c.deletions <- key
		return nil, false
	}
	return &e.value, true
}

// Add adding a new value to the cache with default TTL
func (c *Cache[K, V]) Add(k K, v V) {
	c.lock.Lock()
	defer c.lock.Unlock()
	exp := time.Time{}
	if c.ttl > 0 {
		exp = time.Now().Add(c.ttl)
	}
	e := entry[V]{
		value:     v,
		expiresAt: exp,
	}
	c.items[k] = e
}

// Add adding a new value to the cache with default TTL
func (c *Cache[K, V]) AddWithTTL(k K, v V, ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()
	e := entry[V]{
		value:     v,
		expiresAt: time.Now().Add(ttl),
	}
	c.items[k] = e
}

// Has checking if a value is in the map
func (c *Cache[K, V]) Has(k K) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	e, ok := c.items[k]
	if !ok {
		return false
	}
	if c.isEvicted(e) {
		c.deletions <- k
		return false
	}
	return true
}

// Delete delete a single value
func (c *Cache[K, V]) Delete(k K) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.items, k)
}

// Count getting the count of active values
func (c *Cache[K, V]) Count() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	count := 0
	for _, v := range c.items {
		if !c.isEvicted(v) {
			count++
		}
	}
	return count
}

// Purge remove all values from cache, regardless if evicted or not
func (c *Cache[K, V]) Purge() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.items = make(map[K]entry[V])
}

// DeleteEvicted remove all evicted values from cache
func (c *Cache[K, V]) DeleteEvicted() {
	for k, v := range c.items {
		if c.isEvicted(v) {
			c.deletions <- k
		}
	}
}

// Close closes all needed resources, e.g. the deletion queue channel
func (c *Cache[K, V]) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Stop()
	close(c.deletions)
	close(c.doneAutodelete)
}

// deleteEvicted delete an entry only if it's evicted
// needed for the scenario that the go routine for deletion of evicted entries
// is not fast enough, and a new entry with the same key is already written to the cache
func (c *Cache[K, V]) deleteEvicted(k K) {
	c.lock.Lock()
	defer c.lock.Unlock()
	v, ok := c.items[k]
	if !ok {
		// already deleted
		return
	}
	if c.isEvicted(v) {
		delete(c.items, k)
	}
}

// isEvicted checking if an entry is evicted
func (c *Cache[K, V]) isEvicted(e entry[V]) bool {
	// cache in nottl mode
	if c.ttl < 0 {
		return false
	}
	// item has no expire set
	if e.expiresAt.IsZero() {
		return false
	}
	return e.expiresAt.Before(time.Now())
}
