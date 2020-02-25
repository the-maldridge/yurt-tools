package main

import (
	"log"
	"os"
	"path"

	"github.com/the-maldridge/yurt-tools/internal/consul"
	"github.com/the-maldridge/yurt-tools/internal/nomad"
)

func main() {
	prefix := os.Getenv("CONSUL_PREFIX")
	if prefix == "" {
		prefix = "yurt-tools"
	}

	cs, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}

	nc, err := nomad.New()
	if err != nil {
		log.Fatal(err)
	}

	tasks, err := cs.KnownTasks()
	if err != nil {
		log.Fatal(err)
	}

	for _, task := range tasks {
		switch task.Driver {
		case "docker":
			meta := map[string]string{
				"TRIVY_CONTAINER": task.Docker.RepoStr(),
				"TRIVY_REGISTRY":  task.Docker.Registry().Name,
				"TRIVY_OUTPUT": path.Join(
					prefix,
					"taskinfo",
					task.Job,
					task.Group,
					task.Name,
					"trivy",
				),
			}
			log.Println(meta)
			nc.Dispatch("trivy-scan", meta)
		}
	}
}
