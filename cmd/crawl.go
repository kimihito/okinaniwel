package cmd

import (
	"github.com/kimihito/okinaniwel/crawler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(crawlCmd)

	viper.SetDefault("log_level", "debug")
}

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl dogs from pref-okinawa",
	Run: func(cmd *cobra.Command, args []string) {
		crawler.Run()
	},
}
