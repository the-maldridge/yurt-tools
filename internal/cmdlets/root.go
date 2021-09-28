package cmdlets

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "yurt",
		Short: "yurt is a collection of tools for the prepared Nomad admin",
		Long: rootCmdLongDocs,
	}
	rootCmdLongDocs = `yurt is a complete suite of tools that you can use to make your life
easier maintaining a Nomad cluster.  yurt is a multi-call binary that
you can use to invoke many yurt tools.`
)

// Entrypoint is the entrypoint into all cmdlets, it will dispatch to
// the right one.
func Entrypoint() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
