package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/conradludgate/pass/pass"
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
	fmt.Println("Pairing...")

	file, err := os.OpenFile(os.Getenv("PASS_KEY_FILE"), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		file, err = os.OpenFile(".pass.key", os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			fmt.Println("Cannot open file. Reason:", err.Error())
			return
		}
	}
	defer file.Close()

	var priv *rsa.PrivateKey

	// Try to read file
	b, err := ioutil.ReadAll(file)
	if err == nil {
		// If read is successful, try to parse PrivateKey as JSON
		priv = &rsa.PrivateKey{}

		if json.Unmarshal(b, priv) != nil ||
			priv.Validate() != nil {
			// If error parsing JSON, or the key is invalid, reject key
			priv = nil
		}
	}

	// Open checksum file
	chksumf, err := os.OpenFile(file.Name()+".sum", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		fmt.Println("Cannot open Checksum file.", err.Error())
		return
	}

	if priv != nil {
		sum1 := sha512.Sum512(b)
		sum2, err := ioutil.ReadAll(chksumf)
		// If the checksum can't be read, reject key
		if err != nil || len(sum2) != sha512.Size {
			priv = nil
		}

		// If checksum isn't equal, reject key
		for i, v := range sum2 {
			if sum1[i] != v {
				priv = nil
				break
			}
		}
	}

	// If key is invalid, remake
	if priv == nil {
		// Create new key
		priv, err = rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			fmt.Println("Error generating key.", err.Error())
			return
		}

		// Format key
		b, err = json.Marshal(priv)
		if err != nil {
			fmt.Println("Error formatting private key.", err.Error())
			return
		}

		// Calculate checksum
		sum := sha512.Sum512(b)

		// Write key and checksum to files
		if _, err = file.Write(b); err != nil {
			fmt.Println("Error writing to key file.", err.Error())
			return
		}
		if _, err = chksumf.Write(sum[:]); err != nil {
			fmt.Println("Error writing to checksum file.", err.Error())
			return
		}
	}

	// Pair
	err = pass.Pair(priv)
	if err != nil {
		panic(err)
	}
}
