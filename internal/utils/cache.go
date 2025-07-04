package utils

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrItemNotFound    = errors.New("item not found")
	ErrIndexOutOfBound = errors.New("index out of bound")
)

type CacheItem[T any] struct {
	Id   string
	Item T
}

type CacheList[T any] struct {
	mu    sync.RWMutex
	items []CacheItem[T]
}

func NewCacheList[T any]() *CacheList[T] {
	return &CacheList[T]{
		items: make([]CacheItem[T], 0),
	}
}

func (c *CacheList[T]) Items() []CacheItem[T] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.items
}

func (c *CacheList[T]) ItemValues() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return MapSlice(c.items, func(item CacheItem[T]) T {
		return item.Item
	})
}

func (c *CacheList[T]) Add(item CacheItem[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item.Id == "" {
		item.Id = uuid.NewString()
	}
	c.items = append(c.items, item)
}

func (c *CacheList[T]) AddAll(item []CacheItem[T]) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = append(c.items, item...)
}

func (c *CacheList[T]) Remove(item CacheItem[T]) error {
	return c.RemoveId(item.Id)
}

func (c *CacheList[T]) RemoveId(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, idx, err := c.getByIdLocked(id)
	if err != nil {
		return err
	}
	return c.removeIndexLocked(idx)
}

func (c *CacheList[T]) removeIndexLocked(index int) error {
	if index < 0 || index >= len(c.items) {
		return ErrIndexOutOfBound
	}
	// Entfernen per Slicing (Reihenfolge wird beibehalten)
	c.items = append(c.items[:index], c.items[index+1:]...)
	return nil
}

func (c *CacheList[T]) GetById(id string) (CacheItem[T], int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.getByIdLocked(id)
}

func (c *CacheList[T]) getByIdLocked(id string) (CacheItem[T], int, error) {
	for idx, item := range c.items {
		if item.Id == id {
			return item, idx, nil
		}
	}
	var zero CacheItem[T]
	return zero, -1, ErrItemNotFound
}

func (c *CacheList[T]) Set(item CacheItem[T]) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item.Id == "" {
		item.Id = uuid.NewString()
		c.items = append(c.items, item)
		return nil
	}
	_, idx, err := c.getByIdLocked(item.Id)
	if err != nil {
		// Falls das Item noch nicht existiert, wird es hinzugefügt.
		if errors.Is(err, ErrItemNotFound) {
			c.items = append(c.items, item)
			return nil
		}
		return err
	}
	c.items[idx] = item
	return nil
}

func (c *CacheList[T]) ContainsId(id string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, _, err := c.getByIdLocked(id)
	return err == nil
}

func (c *CacheList[T]) IsEmpty() bool {
	return len(c.items) == 0
}

func (c *CacheList[T]) Clear() {
	c.items = make([]CacheItem[T], 0)
}
