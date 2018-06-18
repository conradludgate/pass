package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/conradludgate/pass/libpass"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pass",
	Short: "Phone authenticated password manager",

	Run: pair,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func pair(cmd *cobra.Command, args []string) {

	file, err := os.OpenFile(os.Getenv("PASS_KEY_FILE"), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		file, err = os.OpenFile(".pass.key", os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			fmt.Println("Cannot open file. Reason:", err.Error())
			return
		}
	}
	defer file.Close()

	key, err := libpass.NewKeyFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	server := os.Getenv("PASS_SERVER")
	u := libpass.DefaultServer
	if server != "" {
		u, err = url.Parse(server)
		if err != nil {
			u = libpass.DefaultServer
		}
	}

	// Pair
	err = libpass.Pair(key, *u)
	if err != nil {
		fmt.Println("Error occured while pairing.", err.Error())
	}
}
