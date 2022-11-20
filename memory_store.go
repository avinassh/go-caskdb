package caskdb

type MemoryStore struct {
	data map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{make(map[string]string)}
}

func (m *MemoryStore) Get(key string) string {
	return m.data[key]
}

func (m *MemoryStore) Set(key string, value string) {
	m.data[key] = value
}

func (m *MemoryStore) Close() bool {
	return true
}
