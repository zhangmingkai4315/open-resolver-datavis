package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/zhangmingkai4315/open-resolver-datavis/config"
)

func init() {
	serverCmd.PersistentFlags().String("http.listen", "", "http server url")
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a web server for show dns data",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetAppConfig()
		httpServerURL := cmd.Flag("http.listen").Value.String()
		if httpServerURL == "" {
			httpServerURL = cfg.GlobalConfig.Listen
		}
		log.Printf("Try starting web serve at %s\n", httpServerURL)
		// web.Start(httpServerURL)
	},
}
