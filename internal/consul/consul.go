package consul

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/hashicorp/consul/api"

	"github.com/the-maldridge/yurt-tools/internal/nomad"
)

type Consul struct {
	*api.Client

	prefix string
}

type TaskData map[string]interface{}

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

func (c *Consul) UpdateTaskData(t nomad.Task, key string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	pair := api.KVPair{
		Key:   path.Join(c.prefix, "taskinfo", t.Job, t.Group, t.Name, key),
		Value: bytes,
	}

	_, err = c.KV().Put(&pair, nil)
	if err != nil {
		return err
	}
	return nil
}

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

func (c *Consul) LoadAllForTask(t nomad.Task) (TaskData, error) {
	m := make(TaskData)

	kvps, _, err := c.KV().List(path.Join(c.prefix, "taskinfo", t.Job, t.Group, t.Name), nil)
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
