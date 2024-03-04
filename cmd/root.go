package cmd

import (
	"awesome-proxy/config"
	"awesome-proxy/proxy"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	Version   = "unknown"
	BuildTime = "unknown"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use: "awesome-proxy",
	Run: func(cmd *cobra.Command, args []string) {
		proxy.Run(config.AppConfig)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&config.AppConfigPath, "config", "", "config file (default is conf/conf.yaml)")
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("version:", Version, "build time:", BuildTime)
		os.Exit(0)
	},
}

func initConfig() {
	if config.AppConfigPath == "" {
		config.AppConfigPath = "conf/conf.yaml"
	}
	config.Init(config.AppConfigPath)
}
