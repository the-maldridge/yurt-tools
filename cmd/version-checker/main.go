package main

import (
	"log"

	"github.com/the-maldridge/yurt-tools/internal/consul"
	"github.com/the-maldridge/yurt-tools/internal/docker"
	"github.com/the-maldridge/yurt-tools/internal/versions"
)

func main() {
	cs, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}

	ds, err := docker.New()
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
			tags, err := ds.GetTags(task.Docker)
			if err != nil {
				log.Printf("Could not load tasks for %s: %v", task.Name, err)
				continue
			}

			info := versions.Compare(task.Docker.Tag, tags)
			if err := cs.UpdateTaskData(task, "versions", info); err != nil {
				log.Printf("Could not update task data for %s: %v", task.Name, err)
			}
		}
	}
}
