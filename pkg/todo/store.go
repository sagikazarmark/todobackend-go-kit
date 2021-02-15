package todo

import (
	"context"
	"sort"
	"sync"
)

// InMemoryStore keeps items in the memory.
// Use it in tests or for development/demo purposes.
type InMemoryStore struct {
	items     map[string]Item
	itemsOnce sync.Once
	mu        sync.RWMutex
}

// NewInMemoryStore returns a new in-memory item store.
func NewInMemoryStore() *InMemoryStore {
	store := &InMemoryStore{}

	store.init()

	return store
}

func (s *InMemoryStore) init() {
	s.itemsOnce.Do(func() {
		s.items = make(map[string]Item)
	})
}

// Store stores an item.
func (s *InMemoryStore) Store(_ context.Context, item Item) error {
	s.init()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[item.ID] = item

	return nil
}

// GetAll returns all items.
func (s *InMemoryStore) GetAll(_ context.Context) ([]Item, error) {
	s.init()

	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]Item, len(s.items))

	// This makes sure items are always returned in the same, sorted order
	keys := make([]string, 0, len(s.items))
	for k := range s.items {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, key := range keys {
		items[i] = s.items[key]
	}

	return items, nil
}

// DeleteItems deletes all items from the store.
func (s *InMemoryStore) DeleteAll(_ context.Context) error {
	s.init()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = make(map[string]Item)

	return nil
}

// GetOne returns a single item by its ID.
func (s *InMemoryStore) GetOne(_ context.Context, id string) (Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	if !ok {
		return item, NotFoundError{ID: id}
	}

	return item, nil
}

// DeleteOne deletes a single item by its ID.
func (s *InMemoryStore) DeleteOne(_ context.Context, id string) error {
	s.init()

	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.items[id]
	if !ok {
		return nil
	}

	delete(s.items, id)

	return nil
}
