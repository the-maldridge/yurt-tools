package cmdlets

import (
	"github.com/spf13/cobra"
)

var (
	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "http provides http servelets for various features and functions",
		Long:  httpCmdLongDocs,
	}
	httpCmdLongDocs = `http cmdlets provide various servers for different parts of the
yurt-tools stack.`
)

func init() {
	rootCmd.AddCommand(httpCmd)
}
