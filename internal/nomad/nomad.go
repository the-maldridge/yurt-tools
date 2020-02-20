package nomad

import (
	"github.com/hashicorp/nomad/api"

	"github.com/the-maldridge/yurt-tools/internal/docker"
)

// A Client is a wrapper around a nomad API client which allows much
// more convenient methods to be built on top of it to provide
// information to the caller and simplify task dispatch.
type Client struct {
	*api.Client
}

type QueryOpts struct {
	Region    string
	Namespace string
	Prefix    string
}

type Task struct {
	Job    string
	Name   string
	Driver string
	URL    string
	Docker docker.Image
}

func New() (*Client, error) {
	c, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	x := new(Client)
	x.Client = c
	return x, nil
}

// ListTasks crawls nomad for running tasks.  This can include batch
// tasks, service tasks, system tasks, and any other kind of task as
// long as it is running.
func (c *Client) ListTasks(cfg QueryOpts) ([]Task, error) {
	qopts := api.QueryOptions{
		Region:    cfg.Region,
		Namespace: cfg.Namespace,
		Prefix:    cfg.Prefix,
	}

	list, _, err := c.Jobs().List(&qopts)
	if err != nil {
		return nil, err
	}

	tl := []Task{}
	for _, i := range list {
		if i.Stop || i.Type == api.JobTypeBatch {
			continue
		}

		job, _, err := c.Jobs().Info(i.ID, nil)
		if err != nil {
			return nil, err
		}
		for _, taskGroup := range job.TaskGroups {
			for _, task := range taskGroup.Tasks {
				t := Task{
					Job:    *job.Name,
					Name:   task.Name,
					Driver: task.Driver,
				}

				switch t.Driver {
				case "docker":
					t.Docker = docker.ParseIdentifier(task.Config["image"].(string))
					t.URL = t.Docker.URL()
				}

				tl = append(tl, t)
			}
		}
	}
	return tl, nil
}
