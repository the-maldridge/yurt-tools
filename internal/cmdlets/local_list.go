package cmdlets

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/the-maldridge/yurt-tools/internal/nomad"
)

var (
	localListCmd = &cobra.Command{
		Use:   "list",
		Short: "list resources available in Nomad",
		Long:  localListCmdLongDocs,
		Run:   localListCmdRun,
	}
	localListCmdLongDocs = `list can be used to see what resources are visible to the yurt-tools.
Its intended to be run from a local workstation to be able to debug
token access issues, and it does not consult with the consul key/value
store.`
)

func init() {
	localCmd.AddCommand(localListCmd)
}

func localListCmdRun(c *cobra.Command, args []string) {
	nc, err := nomad.New()
	if err != nil {
		log.Fatal(err)
	}

	namespaces, err := nc.ListNamespaces()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("The following namespaces will be queried: %s", namespaces)

	for _, n := range namespaces {
		tasks, err := nc.ListTasks(nomad.QueryOpts{Namespace: n})
		if err != nil {
			log.Printf("Error querying namespace %s: %v", n, err)
			continue
		}
		log.Printf("There are %d tasks in namespace %s", len(tasks), n)
	}
}
