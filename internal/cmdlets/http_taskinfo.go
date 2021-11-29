package cmdlets

import (
	"github.com/spf13/cobra"
	"log"

	"github.com/the-maldridge/yurt-tools/internal/http"
	"github.com/the-maldridge/yurt-tools/internal/http/taskinfo"
)

var (
	httpTaskInfoCmd = &cobra.Command{
		Use:   "taskinfo",
		Short: "taskinfo provides a web server that serves JSON information for each task",
		Long:  httpTaskInfoCmdLongDocs,
		Run:   httpTaskInfoCmdRun,
	}
	httpTaskInfoCmdLongDocs = `taskinfo provides 2 endpoints that may be queried:
/all - Provides all tasks in one large array
/detail/{namespace}/{job}/{group}/{task} - Detail for one task

These two endpoints will always return JSON encoded data.
`
)

func init() {
	httpCmd.AddCommand(httpTaskInfoCmd)
}

func httpTaskInfoCmdRun(c *cobra.Command, args []string) {
	srv, err := http.New()
	if err != nil {
		log.Printf("Could not initialize webserver: %s", err)
		return
	}

	h, err := taskinfo.New()
	if err != nil {
		log.Printf("Could not intialize component: %s", err)
		return
	}

	srv.Mount("/taskinfo", h.HTTPEntry())

	srv.Serve(":8080")
}
