package cmd

import (
	"log/slog"

	"vcrawler/internal/crawler"
	"vcrawler/internal/stores/adidas"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the crawler",
	Long:  `Starts the crawler to crawl the store data`,
	Run: func(cmd *cobra.Command, args []string) {
		crawler := crawler.GetCrawler()

		slog.Info("Starting api crawler")
		if err := crawler.Start(adidas.GetAdidasStore()); err != nil {
			slog.Error("Error at starting api crawler", "cause", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
