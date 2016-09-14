package cmd

// extern void systemFork();
import "C"

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yroffin/jarvis-go-ext/server"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the jarvis module",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetString("daemon") == "true" {
			C.systemFork()
		}
		server.Start()
	},
}

func init() {
	viper.AutomaticEnv()

	startFlags := startCmd.Flags()

	startFlags.Int("jarvis.module.port", 7000, "set the listening jarvis module port")
	startFlags.String("jarvis.module.interface", "0.0.0.0", "set the listening jarvis module interface")
	startFlags.String("jarvis.module.name", "module", "set the listening jarvis module name")
	startFlags.String("jarvis.server.url", "http://0.0.0.0:8082", "set the listening jarvis server url")
	startFlags.String("daemon", "true", "set daemon mode")

	viper.BindPFlag("jarvis.module.port", startFlags.Lookup("jarvis.module.port"))
	viper.BindPFlag("jarvis.module.name", startFlags.Lookup("jarvis.module.name"))
	viper.BindPFlag("jarvis.module.interface", startFlags.Lookup("jarvis.module.interface"))
	viper.BindPFlag("jarvis.server.url", startFlags.Lookup("jarvis.server.url"))
	viper.BindPFlag("daemon", startFlags.Lookup("daemon"))

	RootCmd.AddCommand(startCmd)
}
