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

func (c *Consul) UpdateTaskData(t nomad.Task) error {
	bytes, err := json.Marshal(t)
	if err != nil {
		return err
	}

	pair := api.KVPair{
		Key:   path.Join(c.prefix, "taskinfo", t.Job, t.Group, t.Name, "metadata"),
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
		}
		out = append(out, t)
	}
	return out, nil
}
