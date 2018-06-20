package libpass

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server  string
	Keyfile string
}

func (c *Config) LoadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/pass")
	viper.AddConfigPath("$HOME/.pass")

	if err := viper.Unmarshal(c); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
