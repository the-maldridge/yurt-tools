package cmdlets

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/the-maldridge/yurt-tools/internal/kv"
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

	cs, err := kv.NewKVBackend()
	if err != nil {
		log.Fatal(err)
	}

	namespaces, err := nc.ListNamespaces()
	if err != nil {
		log.Fatal(err)
	}

	tasks := make(map[string]nomad.Task)
	for _, n := range namespaces {
		tl, err := nc.ListTasks(nomad.QueryOpts{Namespace: n})
		if err != nil {
			log.Printf("Error querying namespace %s: %v", n, err)
			continue
		}
		for _, t := range tl {
			tasks[t.Path()] = t
		}
	}

	kt, err := cs.KnownTasks()
	if err != nil {
		log.Printf("Error reading back task list: %v", err)
		return
	}
	known := make(map[string]nomad.Task, len(kt))
	for _, k := range kt {
		known[k.Path()] = k
	}

	for path, task := range tasks {
		if err := cs.UpdateTaskData(task, "metadata", task); err != nil {
			log.Printf("Could not update task metadata: %v", err)
		}
		delete(known, path)
	}

	for k, t := range known {
		log.Printf("Obsolete Task: %s", k)
		if err := cs.DeleteTask(t); err != nil {
			log.Printf("Could not remove obsolete task: %v", err)
		}
	}
}
