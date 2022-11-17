package kv

import (
	"os"

	"github.com/hashicorp/consul/api"
)

// Consul provides durable key-value storage on top of a key/value
// store, likely the one implementing service discovery for the nomad
// cluster being scanned.
type Consul struct {
	*api.Client

	prefix string
}

// NewConsul sets up a new consul client and returns a Store
func NewConsul() (Store, error) {
	c, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	x := &Consul{Client: c}

	x.prefix = "yurt-tools"
	if prefix := os.Getenv("CONSUL_PREFIX"); prefix != "" {
		x.prefix = prefix
	}

	return x, nil
}

// Get returns the value at the given key, or an error if no such
// value exists
func (c *Consul) Get(key string) ([]byte, error) {
	kvp, _, err := c.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	return kvp.Value, nil
}

// Put stores a value to the given key
func (c *Consul) Put(key string, val []byte) error {
	pair := api.KVPair{
		Key:   key,
		Value: val,
	}
	_, err := c.KV().Put(&pair, nil)
	return err
}

// ListPrefix takes a prefix and lists all keys at or below that
// prefix.
func (c *Consul) ListPrefix(prefix string) ([]string, error) {
	keys, _, err := c.KV().Keys(prefix, "", nil)
	return keys, err
}

// DelPrefix clears an entire subtree rooted at the given prefix.
func (c *Consul) DelPrefix(prefix string) error {
	_, err := c.KV().DeleteTree(prefix, nil)
	return err
}
