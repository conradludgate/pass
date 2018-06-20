package cmd

import (
	"fmt"
	"os"

	"github.com/conradludgate/pass/libpass"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "pass",
	Long: "Pass is a program designed to securely store passwords on your phone but still being able to use them on your computer. Visit https://pass.conradludgate.com for more details.",
}

var config libpass.Config

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.OnInitialize((&config).LoadConfig)

	rootCmd.PersistentFlags().StringVarP(&config.Keyfile, "keyfile", "k", ".pass.key", "The location of the keyfile")
	rootCmd.PersistentFlags().StringVarP(&config.Server, "server", "s", libpass.DefaultServer, "The pass server to connect to")

	viper.BindPFlag("keyfile", rootCmd.PersistentFlags().Lookup("keyfile"))
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
