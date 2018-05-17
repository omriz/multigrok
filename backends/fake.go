package backends

import (
	"context"
	"fmt"
)

type FakeBackend struct {
	Id string
}

func (b *FakeBackend) Query(ctx context.Context, q string) (WebServiceResult, error) {
	return WebServiceResult{}, fmt.Errorf("Unimplemnted")
}
func (b *FakeBackend) Fetch(ctx context.Context, prefix, path string) ([]byte, error) {
	// Simple implementation for testing.
	return []byte(b.Id), nil
}
func (b *FakeBackend) UID() string { return b.Id }
