package backends

import "fmt"

type FakeBackend struct {
	Id string
}

func (b *FakeBackend) Query(q string) (WebServiceResult, error) {
	return WebServiceResult{}, fmt.Errorf("Unimplemnted")
}
func (b *FakeBackend) Fetch(prefix, path string) ([]byte, error) {
	// Simple implementation for testing.
	return []byte(b.Id), nil
}
func (b *FakeBackend) UID() string { return b.Id }
