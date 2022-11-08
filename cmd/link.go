package cmd

import (
	"fmt"
	"os"

	"github.com/hugowangler/sesame/internal/config"
	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link [URL] [repo] [alias]", // repo optional if inside dir, alias always optional. We should use flags instead
	Short: "Associates a link to a repository",
	Long: `Associates a link to a repository

The associated links will have an index and an optional alias if one is
provided. If you are inside the repository that you want to add a link to you
don't have to specify the repository name.
`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args =", args)
		fmt.Printf("adding a new link to the repository: repo=%s, URL=%s\n", args[1], args[0])
		index, err := config.StoreLink(args[1], args[0])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "unable to store link for the repository: repo=%s err=%v\n", args[1], err)
			os.Exit(1)
		}
		fmt.Printf("added %d new repositories to Sesame\n", index)
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
