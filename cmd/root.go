package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sesame",
	Short: "tool that opens git repositories in your browser",
	Long:  "Sesame is a CLI tool that helps you quickly navigate to your Git projects with ease.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/.sesame.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var cfgPath string
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		cfgPath = home + "/.config/"
		viper.AddConfigPath(cfgPath)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".sesame")
	}

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			_, _ = fmt.Fprintf(os.Stderr, "could not read config file: path=%s, err=%v\n", viper.ConfigFileUsed(), err)
			os.Exit(1)
		}
	}
}
