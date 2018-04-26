package backends

import "fmt"

type FakeBackend struct {
	Id string
}

func (b *FakeBackend) Query(q string) (WebServiceResult, error) {
	return WebServiceResult{}, fmt.Errorf("Unimplemnted")
}
func (b *FakeBackend) Fetch(prefix, path string) ([]byte, error) {
	return nil, fmt.Errorf("Unimplemented")
}
func (b *FakeBackend) UID() string { return b.Id }
