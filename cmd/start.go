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
		server.Start()
	},
}

func fork() {
	C.systemFork()
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start the jarvis module in daemon mode",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		fork()
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

	viper.BindPFlag("jarvis.module.port", startFlags.Lookup("jarvis.module.port"))
	viper.BindPFlag("jarvis.module.name", startFlags.Lookup("jarvis.module.name"))
	viper.BindPFlag("jarvis.module.interface", startFlags.Lookup("jarvis.module.interface"))
	viper.BindPFlag("jarvis.server.url", startFlags.Lookup("jarvis.server.url"))

	RootCmd.AddCommand(startCmd)

	daemonFlags := daemonCmd.Flags()

	daemonFlags.Int("jarvis.module.port", 7000, "set the listening jarvis module port")
	daemonFlags.String("jarvis.module.interface", "0.0.0.0", "set the listening jarvis module interface")
	daemonFlags.String("jarvis.module.name", "module", "set the listening jarvis module name")
	daemonFlags.String("jarvis.server.url", "http://0.0.0.0:8082", "set the listening jarvis server url")

	viper.BindPFlag("jarvis.module.port", daemonFlags.Lookup("jarvis.module.port"))
	viper.BindPFlag("jarvis.module.name", daemonFlags.Lookup("jarvis.module.name"))
	viper.BindPFlag("jarvis.module.interface", daemonFlags.Lookup("jarvis.module.interface"))
	viper.BindPFlag("jarvis.server.url", daemonFlags.Lookup("jarvis.server.url"))

	RootCmd.AddCommand(daemonCmd)
}
