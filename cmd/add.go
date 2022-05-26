package cmd

import (
	"fmt"
	"github.com/hugowangler/sesame/internal/config"
	"github.com/hugowangler/sesame/internal/git"
	"github.com/spf13/cobra"
	"os"
)

var addCmd = &cobra.Command{
	Use:   "add [PATH]",
	Short: "Adds any found repositories starting from PATH",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("searching for new repositories in %s\n", args[0])
		repos, err := git.FindRepos(args[0])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not find repositories: path=%s, err=%v\n", args[0], err)
			os.Exit(1)
		}
		numStored, err := config.StoreRepositories(repos)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not store found repositories in config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("added %d new repositories to Sesame\n", numStored)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
