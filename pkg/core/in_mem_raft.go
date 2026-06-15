package core

import (
	"maps"

	"github.com/mateenbagheri/memorabilia/pkg/types"
)

// Dump returns a shallow copy of the store for Raft snapshotting.
// The copy is taken under a read lock so ongoing reads are not blocked.
func (imc *InMemoryCommandRepository) Dump() (map[string]types.ColumnValueWithTTL, error) {
	imc.mu.RLock()
	defer imc.mu.RUnlock()

	dst := make(map[string]types.ColumnValueWithTTL, len(imc.store))
	src := imc.store
	maps.Copy(dst, src)

	return dst, nil
}

// Load atomically replaces the store with the provided data.
// Used by FSM.Restore() when a follower catches up from a snapshot.
func (imc *InMemoryCommandRepository) Load(src map[string]types.ColumnValueWithTTL) error {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	dst := make(map[string]types.ColumnValueWithTTL, len(src))
	maps.Copy(dst, src)

	imc.store = dst

	return nil
}
