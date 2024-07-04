package core

import "sync"

var _ CommandRepository = (*InMemoryCommandRepository)(nil)

// InMemoryCommandRepository is an in-memory implementation of CommandRepository.
type InMemoryCommandRepository struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewInMemoryCommandRepository() *InMemoryCommandRepository {
	return &InMemoryCommandRepository{
		store: make(map[string]string),
	}
}
