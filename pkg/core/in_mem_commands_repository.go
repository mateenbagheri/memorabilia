package core

import (
	"sync"

	"github.com/mateenbagheri/memorabilia/pkg/types"
)

var _ CommandsRepository = (*InMemoryCommandRepository)(nil)

// InMemoryCommandRepository is an in-memory implementation of CommandRepository.
type InMemoryCommandRepository struct {
	mu    sync.RWMutex
	store map[string]types.ColumnValueWithTTL
}

func NewInMemoryCommandRepository() *InMemoryCommandRepository {
	return &InMemoryCommandRepository{
		store: make(map[string]types.ColumnValueWithTTL),
	}
}

func NewInMemoryCommandRepositoryWithInitialStore(store map[string]types.ColumnValueWithTTL) *InMemoryCommandRepository {
	return &InMemoryCommandRepository{
		store: store,
	}
}
