package cmd

import (
	"fmt"
	"github.com/hugowangler/sesame/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
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

var removeCmd = &cobra.Command{
	Use:   "remove [REPO NAME(s)]",
	Short: "Removes one or more entries from the config file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.GetConfig()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not get config: %v\n", err)
			os.Exit(1)
		}
		pre := len(conf.Repositories)
		for _, arg := range args {
			delete(conf.Repositories, strings.ToLower(arg))
		}
		viper.Set("repositories", conf.Repositories)
		err = viper.WriteConfig()
		if err != nil {
			_, _ = fmt.Fprintf(
				os.Stderr,
				"could not write config file: path=%s, err=%v\n",
				viper.ConfigFileUsed(),
				err,
			)
			os.Exit(1)
		}
		fmt.Printf("removed %d entries from the config file\n", pre-len(conf.Repositories))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(viewCmd)
	configCmd.AddCommand(removeCmd)
}
