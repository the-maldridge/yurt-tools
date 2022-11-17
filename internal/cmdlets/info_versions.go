package cmdlets

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/the-maldridge/yurt-tools/internal/docker"
	"github.com/the-maldridge/yurt-tools/internal/kv"
	"github.com/the-maldridge/yurt-tools/internal/versions"
)

var (
	infoVersionsCmd = &cobra.Command{
		Use:   "version-check",
		Short: "Seek newer versions of task artifacts",
		Long:  versionsCmdLongDocs,
		Run:   versionsCmdRun,
	}
	versionsCmdLongDocs = `The version checker reaches out to remote repositories to determine
newer versions of available containers.  Versions will be compared
using go-version to determine up to 5 versions newer than what is
currently deployed.`
)

func init() {
	infoCmd.AddCommand(infoVersionsCmd)
}

func versionsCmdRun(c *cobra.Command, args []string) {
	cs, err := kv.NewKVBackend()
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
