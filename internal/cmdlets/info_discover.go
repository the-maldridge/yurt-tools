package cmdlets

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/the-maldridge/yurt-tools/internal/consul"
	"github.com/the-maldridge/yurt-tools/internal/nomad"
)

var (
	infoDiscoverCmd = &cobra.Command{
		Use:   "discover",
		Short: "discover tasks running in Nomad",
		Long:  discoverCmdLongDocs,
		Run:   discoverCmdRun,
	}
	discoverCmdLongDocs = `scrape a listing of all tasks running in Nomad and store them for
other tasks to key off of.`
)

func init() {
	infoCmd.AddCommand(infoDiscoverCmd)
}

func discoverCmdRun(c *cobra.Command, args []string) {
	nc, err := nomad.New()
	if err != nil {
		log.Fatal(err)
	}

	cs, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}

	namespaces, err := nc.ListNamespaces()
	if err != nil {
		log.Fatal(err)
	}

	tasks := []nomad.Task{}
	for _, n := range namespaces {
		t, err := nc.ListTasks(nomad.QueryOpts{Namespace: n})
		if err != nil {
			log.Printf("Error querying namespace %s: %v", n, err)
			continue
		}
		tasks = append(tasks, t...)
	}

	for _, task := range tasks {
		if err := cs.UpdateTaskData(task, "metadata", task); err != nil {
			log.Printf("Could not update task metadata: %v", err)
		}
	}
}
