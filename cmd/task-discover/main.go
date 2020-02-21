package main

import (
	"log"

	"github.com/the-maldridge/yurt-tools/internal/consul"
	"github.com/the-maldridge/yurt-tools/internal/nomad"
)

func main() {
	nc, err := nomad.New()
	if err != nil {
		log.Fatal(err)
	}

	cs, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}

	tasks, err := nc.ListTasks(nomad.QueryOpts{})
	if err != nil {
		log.Fatal(err)
	}

	for _, task := range tasks {
		if err := cs.UpdateTaskData(task, "metadata", task); err != nil {
			log.Printf("Could not update task metadata: %v", err)
		}
	}
}
