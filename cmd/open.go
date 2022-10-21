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
	Run: func(cmd *cobra.Command, args []string) {
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
