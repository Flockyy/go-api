package storage

import (
	"sync"
	"time"

	"go-api/models"

	"github.com/google/uuid"
)

// Store interface defines the contract for data storage
type Store[T any] interface {
	GetAll() []T
	GetByID(id string) (T, bool)
	Create(data T) T
	Update(id string, data T) (T, bool)
	Delete(id string) bool
}

// MemoryStore implements Store interface with in-memory storage
type MemoryStore[T any] struct {
	mu    sync.RWMutex
	items map[string]T
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore[T any]() *MemoryStore[T] {
	return &MemoryStore[T]{
		items: make(map[string]T),
	}
}

// GetAll returns all items
func (s *MemoryStore[T]) GetAll() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]T, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}
	return items
}

// GetByID retrieves an item by ID
func (s *MemoryStore[T]) GetByID(id string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[id]
	return item, exists
}

// Create adds a new item
func (s *MemoryStore[T]) Create(data T) T {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Use reflection to set ID and timestamps for models
	switch v := any(&data).(type) {
	case *models.Item:
		v.ID = uuid.New().String()
		v.CreatedAt = time.Now()
		v.UpdatedAt = time.Now()
		s.items[v.ID] = any(*v).(T)
	case *models.Client:
		v.ID = uuid.New().String()
		v.CreatedAt = time.Now()
		v.UpdatedAt = time.Now()
		s.items[v.ID] = any(*v).(T)
	}

	return data
}

// Update modifies an existing item
func (s *MemoryStore[T]) Update(id string, data T) (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[id]; !exists {
		var zero T
		return zero, false
	}

	// Preserve ID and CreatedAt, update UpdatedAt
	switch v := any(&data).(type) {
	case *models.Item:
		old := s.items[id]
		oldItem := any(old).(models.Item)
		v.ID = id
		v.CreatedAt = oldItem.CreatedAt
		v.UpdatedAt = time.Now()
		s.items[id] = any(*v).(T)
	case *models.Client:
		old := s.items[id]
		oldClient := any(old).(models.Client)
		v.ID = id
		v.CreatedAt = oldClient.CreatedAt
		v.UpdatedAt = time.Now()
		s.items[id] = any(*v).(T)
	}

	return data, true
}

// Delete removes an item
func (s *MemoryStore[T]) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[id]; !exists {
		return false
	}

	delete(s.items, id)
	return true
}
