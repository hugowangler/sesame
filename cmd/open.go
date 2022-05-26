package cmd

import (
	"fmt"
	"github.com/hugowangler/sesame/internal/config"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"os"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open [REPO NAME]",
	Short: "Opens a repository that is stored in your config in your browser",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("open called")
		config, err := config.GetConfig()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not open repository: %v\n", err)
			os.Exit(1)
		}
		url, exists := config.Repositories[args[0]]
		if !exists {
			_, _ = fmt.Fprintf(os.Stderr, "unknown repository, please add it first: name=%s\n", args[0])
			os.Exit(1)
		}
		err = browser.OpenURL(url)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not open repository: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
