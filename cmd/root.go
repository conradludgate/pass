package cmd

import (
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/conradludgate/pass/libpass"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ed25519"
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

	key, sum, err := ReadKey(file)

	// If key is invalid, remake
	if err != nil {
		fmt.Println(err)
		// Create new key
		_, sec, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			fmt.Println("Error generating key.", err.Error())
			return
		}

		// Calculate checksum
		sum = sha512.Sum512(sec)

		// Write key and checksum to files

		file.Truncate(0)
		file.Seek(0, 0)

		if _, err = file.Write(append(sec, sum[:]...)); err != nil {
			fmt.Println("Error writing to file.", err.Error())
			return
		}

		copy(key[:], sec)
	}

	server := os.Getenv("PASS_SERVER")
	if server != nil {
		u, err := url.Parse(server)
		if err != nil {
			pass.PASS_SERVER = u
		}
	}

	// Pair
	err = pass.Pair(key)
	if err != nil {
		panic(err)
	}
}

func ReadKey(r io.Reader) (key, sum [64]byte, err error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	fmt.Println(len(b))

	if len(b) != 128 {
		return key, sum, errors.New("File must contain only 128 bytes")
	}

	copy(key[:], b)
	copy(sum[:], b[64:])

	sum = sha512.Sum512(b[:64])
	for i, v := range sum {
		if b[64+i] != v {
			return key, sum, errors.New("Bad checksum")
		}
	}

	copy(key[:], b[:64])
	copy(sum[:], b[64:])
	return
}
