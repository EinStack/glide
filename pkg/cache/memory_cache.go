package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/EinStack/glide/pkg/api/schemas"
)

type CacheEntry struct {
	Response  schemas.ChatResponse
	Timestamp time.Time
}

type MemoryCache struct {
	cache map[string]CacheEntry
	mux   sync.Mutex
}

func NewMemoryCache(ttl int, maxSize int) *MemoryCache {
	return &MemoryCache{
		cache: make(map[string]CacheEntry),
	}
}

func (m *MemoryCache) Get(key string) (schemas.ChatResponse, bool) {
	m.mux.Lock()
	defer m.mux.Unlock()
	entry, exists := m.cache[key]
	if !exists {
		return schemas.ChatResponse{}, false
	}
	return entry.Response, true
}

func (m *MemoryCache) Set(key string, response schemas.ChatResponse) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.cache[key] = CacheEntry{
		Response:  response,
		Timestamp: time.Now(),
	}
}

func (m *MemoryCache) All() {
	m.mux.Lock()
	defer m.mux.Unlock()

	fmt.Println("%v", m.cache)
}
