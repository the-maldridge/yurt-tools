package cmdlets

import (
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/the-maldridge/yurt-tools/internal/consul"
	"github.com/the-maldridge/yurt-tools/internal/nomad"
)

var (
	infoTrivyCmd = &cobra.Command{
		Use:   "trivy-dispatch",
		Short: "Launch trivy scans for all compatible discovered containers",
		Long:  trivyCmdLongDocs,
		Run:   trivyCmdRun,
	}
	trivyCmdLongDocs = `Security scanning of compatible container images is performed by
Aquasec Trivy (https://aquasecurity.github.io/trivy/.  This info
provider can be used to launch highly parallel scans of your
discovered containers to a trivy server.).`
)

func init() {
	infoCmd.AddCommand(infoTrivyCmd)
}

func trivyCmdRun(c *cobra.Command, args []string) {
	prefix := os.Getenv("CONSUL_PREFIX")
	if prefix == "" {
		prefix = "yurt-tools"
	}

	job := os.Getenv("YURT_TRIVY_DISPATCHABLE")
	if job == "" {
		job = "yurt-task-trivy-scan"
	}

	cs, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}

	nc, err := nomad.New()
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
			meta := map[string]string{
				"TRIVY_CONTAINER": task.Docker.RepoStr(),
				"TRIVY_REGISTRY":  task.Docker.Registry().Name,
				"TRIVY_OUTPUT": path.Join(
					prefix,
					"taskinfo",
					task.Job,
					task.Group,
					task.Name,
					"trivy",
				),
			}
			if err := nc.Dispatch(job, meta); err != nil {
				log.Printf("Error dispatching job: %v", err)
			}
		}
	}
}
