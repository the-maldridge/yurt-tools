package kv

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/the-maldridge/yurt-tools/internal/nomad"
)

// NewKVBackend provides the KV interface on top of the generic Store
// interface.
func NewKVBackend() (*KV, error) {
	backend := "consul"
	if b := os.Getenv("YURT_BACKEND"); b != "" {
		backend = b
	}

	var store Store
	var err error
	switch backend {
	case "consul":
		store, err = NewConsul()
	case "redis":
		store, err = NewRedis()
	}

	if err != nil {
		return nil, err
	}

	return &KV{store}, nil
}

// UpdateTaskData creates or updates the data for the specified task.
func (k *KV) UpdateTaskData(t nomad.Task, key string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return k.s.Put(path.Join("taskinfo", t.Namespace, t.Job, t.Group, t.Name, key), bytes)
}

// KnownTasks returns a list of all tasks that are known to
// yurt-tools.  This is not necessarily every task known to nomad, but
// all the ones that yurt-tools was allowed to read.
func (k *KV) KnownTasks() ([]nomad.Task, error) {
	keys, err := k.s.ListPrefix("taskinfo")
	if err != nil {
		return nil, err
	}

	out := []nomad.Task{}
	for _, key := range keys {
		bytes, err := k.s.Get(key)
		if err != nil {
			log.Printf("Error with fetch for key: %s %v", key, err)
			continue
		}

		t := nomad.Task{}
		if err := json.Unmarshal(bytes, &t); err != nil {
			log.Printf("Bad json at key: %s %v", key, err)
			continue
		}
		if t.Name == "" {
			// Empty task
			continue
		}
		out = append(out, t)
	}
	return out, nil
}

// LoadAllForTask loads up all the information for a task including
// the default taskinfo blob, and anything written by other discovery
// plugins.
func (k *KV) LoadAllForTask(t nomad.Task) (TaskData, error) {
	m := make(TaskData)

	keys, err := k.s.ListPrefix(path.Clean(path.Join("taskinfo", t.Path())))
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		bytes, err := k.s.Get(key)
		if err != nil {
			log.Printf("Error pulling value at key: %s - %v", key, err)
			continue
		}

		var tmp interface{}
		if err := json.Unmarshal(bytes, &tmp); err != nil {
			log.Printf("Couldn't unmarshal data for %s:%s - %v", t.Name, key, err)
			continue
		}
		m[path.Base(key)] = tmp
	}
	return m, nil
}

// DeleteTask removes the taskinfo blob and all other information for
// a given task and is mainly intended to be a way for removing
// records of tasks that are no longer running.
func (k *KV) DeleteTask(t nomad.Task) error {
	return k.s.DelPrefix(path.Clean(path.Join("taskinfo", t.Path())))
}
