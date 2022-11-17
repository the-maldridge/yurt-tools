package kv

import (
	"context"
	"os"
	"path"

	"github.com/go-redis/redis/v8"
)

// Redis implements storage on top of a redis server.  This should be
// durable, but this is entirely dependent on the configuration of the
// redis server.
type Redis struct {
	rdb *redis.Client
}

// NewRedis returns a Store on top of a remote redis server.
func NewRedis() (Store, error) {
	addr := "localhost:6379"
	if a := os.Getenv("REDIS_ADDR"); a != "" {
		addr = a
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Redis{rdb}, nil
}

// Get returns the available data or a suitable error.
func (r *Redis) Get(key string) ([]byte, error) {
	res := r.rdb.Get(context.Background(), key)
	if res.Err() != nil {
		return nil, res.Err()
	}
	return res.Bytes()
}

// Put stores data to the named key.
func (r *Redis) Put(key string, value []byte) error {
	res := r.rdb.Set(context.Background(), key, value, 0)
	return res.Err()
}

// ListPrefix returns the list of keys below a given prefix
func (r *Redis) ListPrefix(prefix string) ([]string, error) {
	res := r.rdb.Keys(context.Background(), path.Join(prefix, "*"))
	return res.Result()
}

// DelPrefix removes all data below a given prefix
func (r *Redis) DelPrefix(prefix string) error {
	keys, err := r.ListPrefix(prefix)
	if err != nil {
		return err
	}

	res := r.rdb.Del(context.Background(), keys...)
	return res.Err()
}
