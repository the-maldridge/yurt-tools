package kv

// Store defines the mechanisms that a storage interface must provide.
// These do not strictly need to be long term durable, but that is
// preferable.
type Store interface {
	Get(string) ([]byte, error)
	Put(string, []byte) error
	DelPrefix(string) error
	ListPrefix(string) ([]string, error)
}

// KV is a mechanism for storing and returning information within the
// system about tasks.
type KV struct {
	s Store
}

// TaskData is the base type for information concerning a task running
// in the cluster.
type TaskData map[string]interface{}
