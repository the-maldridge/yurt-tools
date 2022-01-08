package consul

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/hashicorp/consul/api"

	"github.com/the-maldridge/yurt-tools/internal/nomad"
)

// Consul wraps the consul client to provide a slimmed down and
// specialized interface.
type Consul struct {
	*api.Client

	prefix string
}

// TaskData is the base type for information concerning a task running
// in the cluster.
type TaskData map[string]interface{}

// New connects to consul in the specified prefix.
func New() (*Consul, error) {
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

// UpdateTaskData creates or updates the data for the specified task.
func (c *Consul) UpdateTaskData(t nomad.Task, key string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	pair := api.KVPair{
		Key:   path.Join(c.prefix, "taskinfo", t.Namespace, t.Job, t.Group, t.Name, key),
		Value: bytes,
	}

	_, err = c.KV().Put(&pair, nil)
	if err != nil {
		return err
	}
	return nil
}

// KnownTasks returns a list of all tasks that are known to
// yurt-tools.  This is not necessarily every task known to nomad, but
// all the ones that yurt-tools was allowed to read.
func (c *Consul) KnownTasks() ([]nomad.Task, error) {
	kvps, _, err := c.KV().List(path.Join(c.prefix, "taskinfo"), nil)
	if err != nil {
		return nil, err
	}

	out := []nomad.Task{}
	for _, kv := range kvps {
		t := nomad.Task{}
		if err := json.Unmarshal(kv.Value, &t); err != nil {
			log.Printf("Bad json at key: %s %v", kv.Key, err)
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
func (c *Consul) LoadAllForTask(t nomad.Task) (TaskData, error) {
	m := make(TaskData)

	kvps, _, err := c.KV().List(path.Clean(path.Join(c.prefix, "taskinfo", t.Path())), nil)
	if err != nil {
		return nil, err
	}

	for _, kv := range kvps {
		var tmp interface{}
		if err := json.Unmarshal(kv.Value, &tmp); err != nil {
			log.Printf("Couldn't unmarshal data for %s:%s - %v", t.Name, kv.Key, err)
			continue
		}
		m[path.Base(kv.Key)] = tmp
	}
	return m, nil
}

// DeleteTask removes the taskinfo blob and all other information for
// a given task and is mainly intended to be a way for removing
// records of tasks that are no longer running.
func (c *Consul) DeleteTask(t nomad.Task) error {
	_, err := c.KV().DeleteTree(path.Clean(path.Join(c.prefix, "taskinfo", t.Path())), nil)
	return err
}
