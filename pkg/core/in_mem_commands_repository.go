package core

import (
	"sync"

	"github.com/mateenbagheri/memorabilia/pkg/types"
)

var _ CommandRepository = (*InMemoryCommandRepository)(nil)

// InMemoryCommandRepository is an in-memory implementation of CommandRepository.
type InMemoryCommandRepository struct {
	mu    sync.RWMutex
	store map[string]types.ColumnValue
}

func NewInMemoryCommandRepository() *InMemoryCommandRepository {
	return &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValue),
	}
}
