package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zhangmingkai4315/open-resolver-datavis/config"
	"github.com/zhangmingkai4315/open-resolver-datavis/db"
)

var cfgFile string
var rootCmd = &cobra.Command{
	Use:   "open-resolver-datavis",
	Short: "open-resolver-datavis is a tools to collect open dns resolver data and using web to show it",
	Long:  `Project source code is available at https://github.com/zhangmingkai4315/open-resolver-datavis`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./config.toml", "config file for read")

}

func initConfig() {
	appConfig, err := config.NewAppConfig(cfgFile)
	if err != nil || appConfig == nil {
		log.Printf("Read config file error:%s\n", err)
		os.Exit(1)
	}

	_, err = db.InitSessionConnect(appConfig.DatabaseConfig.URL)
	if err != nil {
		log.Printf("Connect mongodb server fail : %s\n", err)
	}
	log.Printf("Connect mongodb server success\n")
}

// Execute will execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
