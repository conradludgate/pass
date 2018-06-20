package cmd

import (
	"fmt"
	"os"

	"github.com/conradludgate/pass/libpass"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var pairCmd = &cobra.Command{
	Use:   "pair",
	Short: "Pair with your phone",
	Long:  "Displays a QR code for you to scan on the Pass smartphone app. Go to https://pass.conradludgate.com for more details",

	Run: pair,
}

var timeout int
var force bool

func init() {
	rootCmd.AddCommand(pairCmd)

	pairCmd.Flags().IntVarP(&timeout, "timeout", "t", 60, "Number of seconds before timeout, max 60")
	pairCmd.Flags().BoolVarP(&force, "force", "f", false, "Ignore any warnings while pairing")
}

func pair(cmd *cobra.Command, args []string) {
	file, err := os.OpenFile(config.Keyfile, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		fmt.Println("Cannot open file.", err.Error())
		return
	}
	defer file.Close()

	key, err := libpass.NewKeyFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Pair
	err = libpass.Pair(key, config.Server, timeout, force)
	if err != nil {
		fmt.Println("Error occured while pairing.", err.Error())
	}
}
