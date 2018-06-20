package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [NAME]",
	Short: "Get a password from your paired phone",
	Run:   get,
}

var u bool

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVarP(&u, "url", "u", false, "Search by URL instead of name")
}

func get(cmd *cobra.Command, args []string) {
	fmt.Println(strings.Join(args, "\n"))
}
