package cmd

import (
	"fmt"
	"github.com/hugowangler/sesame/internal/git"
	"github.com/spf13/cobra"
	"os"
)

var recursive bool

var addCmd = &cobra.Command{
	Use:   "add [PATH to repository]",
	Short: "Adds a new repository",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called with args:", args)
		repos, err := git.FindRepos(args[0])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not find repositories: path=%s, err=%v\n", args[0], err)
		}
		for _, r := range repos {
			fmt.Println("repo=", *r)
			fmt.Println("url=", r.Url())
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	// addCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "recursively searches for repositories inside PATH")
}
