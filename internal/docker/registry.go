package docker

import (
	"os"
	"strings"

	"github.com/nokia/docker-registry-client/registry"
)

type Docker struct {
	clients map[string]*registry.Registry
}

func New() (*Docker, error) {
	return &Docker{clients: make(map[string]*registry.Registry)}, nil
}

func (d *Docker) GetTags(i Image) ([]string, error) {
	ri := i.Registry()
	rh, ok := d.clients[ri.Name]
	if !ok {
		// No client for this registry exists yet.
		prefix := strings.ReplaceAll(strings.ToUpper(ri.Name), "-", "_")
		username := os.Getenv(prefix + "_USERNAME")
		password := os.Getenv(prefix + "_PASSWORD")

		c, err := registry.New(ri.URL, username, password)
		if err != nil {
			return nil, err
		}
		rh = c
		d.clients[ri.Name] = c
	}

	repoStr := i.Owner + "/" + i.Image
	if i.Owner == "" && ri.Name == "docker-hub" {
		repoStr = "library/" + i.Image
	}

	tags, err := rh.Tags(repoStr)
	if err != nil {
		return nil, err
	}
	return tags, nil
}
