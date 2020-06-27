package cmd

import (
	"encoding/json"
	"io/ioutil"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		dogs := crawler.Run()
		file, err := json.MarshalIndent(dogs, "", " ")
		if err != nil {
			return err
		}

		err = ioutil.WriteFile("dist/dogs.json", file, 0644)
		if err != nil {
			return err
		}
		return nil
	},
}
