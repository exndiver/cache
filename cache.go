package cache

import (
	"sync"
	"time"
)

//Storage mecanism for caching strings
type Storage interface {
	Get(key string) []byte
	Set(key string, content []byte, duration time.Duration)
}

type Item struct {
	Content    []byte
	Expiration int64
}

//Storage mecanism for caching strings in memory
type Storage struct {
	items map[string]Item
	mu    *sync.RWMutex
}

//NewStorage creates a new in memory storage
func NewStorage() *Storage {
	return &Storage{
		items: make(map[string]Item),
		mu:    &sync.RWMutex{},
	}
}

//Get a cached content by key
func (s Storage) Get(key string) []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item := s.items[key]
	if item.Expired() {
		delete(s.items, key)
		return nil
	}
	return item.Content
}

//Set a cached content by key
func (s Storage) Set(key string, content []byte, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = Item{
		Content:    content,
		Expiration: time.Now().Add(duration).UnixNano(),
	}
}
