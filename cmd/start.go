package cmd

// extern void systemSighup();
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
		C.systemSighup()
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
	startFlags.String("daemon", "false", "set daemon mode")
	startFlags.String("jarvis.option.mfrc522", "false", "mfrc522 init on start")
	startFlags.String("jarvis.option.wiringpi", "false", "wiringpi init on start")
	startFlags.String("jarvis.option.wiringpi.alt", "true", "alternate code")
	startFlags.String("jarvis.option.advertise", "true", "advertise jarvis")
	startFlags.String("jarvis.option.nfctag", "false", "advertise jarvis for nfc tag detection")

	viper.BindPFlag("jarvis.module.port", startFlags.Lookup("jarvis.module.port"))
	viper.BindPFlag("jarvis.module.name", startFlags.Lookup("jarvis.module.name"))
	viper.BindPFlag("jarvis.module.interface", startFlags.Lookup("jarvis.module.interface"))
	viper.BindPFlag("jarvis.server.url", startFlags.Lookup("jarvis.server.url"))
	viper.BindPFlag("jarvis.option.mfrc522", startFlags.Lookup("jarvis.option.mfrc522"))
	viper.BindPFlag("jarvis.option.wiringpi", startFlags.Lookup("jarvis.option.wiringpi"))
	viper.BindPFlag("jarvis.option.wiringpi.alt", startFlags.Lookup("jarvis.option.wiringpi.alt"))
	viper.BindPFlag("jarvis.option.advertise", startFlags.Lookup("jarvis.option.advertise"))
	viper.BindPFlag("jarvis.option.nfctag", startFlags.Lookup("jarvis.option.nfctag"))
	viper.BindPFlag("daemon", startFlags.Lookup("daemon"))

	RootCmd.AddCommand(startCmd)
}
