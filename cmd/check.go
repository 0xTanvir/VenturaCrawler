package cmd

import (
	"log/slog"

	"vcrawler/internal/crawler"
	"vcrawler/internal/stores/adidas"

	"github.com/spf13/cobra"
)

var dump int

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Performs a test crawl of a single store, and dumps the output.",
	Long: `Performs a test crawl of a single store, and dumps the output. 
	No database connection is performed at all.
	Used for testing new and changed store crawlers.`,
	Run: func(cmd *cobra.Command, args []string) {
		crawler := crawler.GetCrawler()

		if err := crawler.Test(dump, adidas.GetAdidasStore()); err != nil {
			slog.Error("Error at checking crawler", "cause", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().IntVarP(&dump, "dump", "d", 0, "dump limit")
}
