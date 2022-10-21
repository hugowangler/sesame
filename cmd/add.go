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
	Long: `Adds any found repositories starting from the specified PATH.

The repositories are found by performing a walk in the file tree rooted at PATH.
Whenever a .git/ directory is found the git config file is parsed in order to
extract the remote origin. The remote origin is then Regex matched to construct 
the URL that is stored in the Sesame config file.
`,
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
