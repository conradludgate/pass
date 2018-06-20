package cmd

import (
	"fmt"
	"os"

	"github.com/conradludgate/pass/libpass"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify that the key exchange wasn't intercepted",
	Run:   verify,
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}

func verify(cmd *cobra.Command, args []string) {
	file, err := os.OpenFile(config.Keyfile, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		fmt.Println("Cannot open file. Reason:", err.Error())
		return
	}
	defer file.Close()

	key, err := libpass.NewKeyFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Pair
	err = libpass.Verify(key)
	if err != nil {
		fmt.Println(err.Error())
	}
}
