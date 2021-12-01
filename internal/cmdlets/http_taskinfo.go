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
		Short: "taskinfo provides a web server with aggregated information",
		Long:  httpTaskInfoCmdLongDocs,
		Run:   httpTaskInfoCmdRun,
	}
	httpTaskInfoCmdLongDocs = `taskinfo provides 2 endpoints that may be queried:

/data/all - Provides all tasks in one large array
/data/detail/{namespace}/{job}/{group}/{task} - Detail for one task

/view/ - HTML overview of status
/view/{namespace}/ - Same as above, per namespace

Endpoints under /data should be expected to return JSON.  Endpoints
under /view will return HTML.`
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
