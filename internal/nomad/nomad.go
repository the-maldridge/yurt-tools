package nomad

import (
	"path"

	"github.com/hashicorp/nomad/api"

	"github.com/the-maldridge/yurt-tools/internal/docker"
)

// A Client is a wrapper around a nomad API client which allows much
// more convenient methods to be built on top of it to provide
// information to the caller and simplify task dispatch.
type Client struct {
	*api.Client
}

// QueryOpts shadows the nomad type of the same name.
type QueryOpts struct {
	Region    string
	Namespace string
	Prefix    string
}

// Task represents a simplified view of a nomad task.
type Task struct {
	Namespace string
	Job       string
	Group     string
	Name      string
	Driver    string
	URL       string
	Docker    docker.Image
}

// Path returns a path like string representation of where the task
// lives in the hierarchy.
func (t Task) Path() string {
	return "/" + path.Join(t.Namespace, t.Job, t.Group, t.Name)
}

// New returns a new nomad client initialized with parameters from the
// environment.
func New() (*Client, error) {
	c, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	x := new(Client)
	x.Client = c
	return x, nil
}

// ListNamespaces returns a list of namespaces in the given cluster.
func (c *Client) ListNamespaces() ([]string, error) {
	list, _, err := c.Namespaces().List(&api.QueryOptions{})
	if err != nil {
		return nil, err
	}

	out := make([]string, len(list))
	for i := range list {
		out[i] = list[i].Name
	}

	return out, nil
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

		job, _, err := c.Jobs().Info(i.ID, &qopts)
		if err != nil {
			return nil, err
		}
		for _, taskGroup := range job.TaskGroups {
			for _, task := range taskGroup.Tasks {
				t := Task{
					Namespace: cfg.Namespace,
					Job:       *job.Name,
					Group:     *taskGroup.Name,
					Name:      task.Name,
					Driver:    task.Driver,
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

// Dispatch can be used to dispatch a parameterized task.
func (c *Client) Dispatch(name string, kv map[string]string) error {
	if _, _, err := c.Jobs().Dispatch(name, kv, nil, nil); err != nil {
		return err
	}
	return nil
}
