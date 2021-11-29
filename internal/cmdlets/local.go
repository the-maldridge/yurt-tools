package cmdlets

import (
	"github.com/spf13/cobra"
)

var (
	localCmd = &cobra.Command{
		Use:   "local",
		Short: "local cmdlets provide debugging assistance locally.",
		Long:  localCmdLongDocs,
	}
	localCmdLongDocs = `Local commands are intended to run from an administrator's machine
and provide some help in debugging and getting various things
working.`
)

func init() {
	rootCmd.AddCommand(localCmd)
}
