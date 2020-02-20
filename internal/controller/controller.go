package controller

import (
	"log"

	"github.com/the-maldridge/yurt-tools/internal/consul"
	"github.com/the-maldridge/yurt-tools/internal/docker"
	"github.com/the-maldridge/yurt-tools/internal/nomad"
)

func New() (*Controller, error) {
	nc, err := nomad.New()
	if err != nil {
		return nil, err
	}

	ds, err := docker.New()
	if err != nil {
		return nil, err
	}

	cs, err := consul.New()
	if err != nil {
		return nil, err
	}

	return &Controller{
		ds: ds,
		nc: nc,
		cs: cs,
	}, nil
}

func (c *Controller) Discover() error {
	tasks, err := c.nc.ListTasks(nomad.QueryOpts{})
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if err := c.cs.UpdateTaskMeta(task); err != nil {
			log.Printf("Could not update task metadata: %v", err)
		}
	}

	return nil
}
