package controller

import (
	"github.com/the-maldridge/yurt-tools/internal/consul"
	"github.com/the-maldridge/yurt-tools/internal/docker"
	"github.com/the-maldridge/yurt-tools/internal/nomad"
)

type Controller struct {
	nc *nomad.Client
	ds *docker.Docker
	cs *consul.Consul
}
