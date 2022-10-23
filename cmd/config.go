package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Operations related to the config file",
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View the contents of the config file",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() == "" {
			fmt.Printf("no config file found, try adding repositories using `sesame add PATH` first\n")
			return
		}
		configContent, err := os.ReadFile(viper.ConfigFileUsed())
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not read config file: path=%s, err=%v\n", viper.ConfigFileUsed(), err)
			os.Exit(1)
		}
		fmt.Printf("viewing config file in %s:\n\n", viper.ConfigFileUsed())
		fmt.Printf(string(configContent))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(viewCmd)
}
