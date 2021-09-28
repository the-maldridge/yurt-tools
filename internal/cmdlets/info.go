package cmdlets

import (
	"github.com/spf13/cobra"
)

var (
	infoCmd = &cobra.Command{
		Use: "info",
		Short: "info cmdlets scrape information from tasks",
		Long: infoCmdLongDocs,
	}
	infoCmdLongDocs = `info cmdlets provide various information services which
either obtain information from tasks directly, or provide
information based on external data sources.`
)

func init() {
	rootCmd.AddCommand(infoCmd)
}
