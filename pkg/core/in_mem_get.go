package core

import "context"

func (imc *InMemoryCommandRepository) Get(ctx context.Context, key string) (value string, err error) {
	imc.mu.Lock() // TODO: remove in the future. this is only here for satisfying linter.
	return "", nil
}
