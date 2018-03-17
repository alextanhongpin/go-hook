package webhook

// Store represents the interface for the different stores available
type Store interface {
	Put(key string, val []byte) error
	Get(key string) ([]byte, error)
	// Set()
	// List()
}

const (
	// InMemory is an enum that represents the in-memory store
	InMemory = 0
	// Consul is an enum that represents the consul store
	Consul = 1
)

// NewStore returns a new store based on the provided options
func NewStore(store int8) Store {
	switch store {
	case Consul:
		return NewConsulStore()
	case InMemory:
		return NewInMemoryStore()
	default:
		return NewInMemoryStore()
	}
}
