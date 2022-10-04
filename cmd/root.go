package cmd

import (
	"fmt"
	"os"

	"github.com/sboysel/ultralist/ultralist"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "ultralist",
	Short: "Ultralist, simple task management for tech folks.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home := ultralist.UserHomeDir()

		viper.AddConfigPath(home)
		viper.AddConfigPath(home + "/.config/ultralist")
		viper.SetConfigName(".ultralist")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error: Can't load config file:", viper.ConfigFileUsed())
		fmt.Println("Run 'ultralist --help' for usage.")
		os.Exit(1)
	}
}
