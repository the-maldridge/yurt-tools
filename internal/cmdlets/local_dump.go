package cmdlets

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/the-maldridge/yurt-tools/internal/kv"
)

var (
	localDumpCmd = &cobra.Command{
		Use:   "dump",
		Short: "dump resources available in kv taskinfo heirarchy",
		Long:  localDumpCmdLongDocs,
		Run:   localDumpCmdRun,
	}
	localDumpCmdLongDocs = ``
)

func init() {
	localCmd.AddCommand(localDumpCmd)
}

func localDumpCmdRun(c *cobra.Command, args []string) {
	cc, err := kv.NewKVBackend()
	if err != nil {
		log.Fatal(err)
	}

	tasks, err := cc.KnownTasks()
	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(tasks)
}
