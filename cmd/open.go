package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hugowangler/sesame/internal/config"
	"github.com/hugowangler/sesame/internal/git"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open [REPO NAME]",
	Short: "Opens a repository that is stored in your config in your browser.",
	Long: `Opens a repository that is stored in your config in your browser.

If REPO NAME is not specified, sesame will attempt to open the repository that
corresponds to the current working directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.GetConfig()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "could not read the config file: %v\n", err)
			os.Exit(1)
		}

		// attempt to open the current directory if no arguments are provided
		if len(args) == 0 {
			workingDir, err := os.Getwd()
			if err != nil {
				_, _ = fmt.Fprintf(
					os.Stderr,
					"unable to open current directory: could not get working directory path: %v\n",
					err,
				)
				os.Exit(1)
			}
			dirSplit := strings.Split(filepath.ToSlash(workingDir), "/")
			dirName := dirSplit[len(dirSplit)-1]
			entry, exists := conf.Repositories[dirName]
			// try to find the repository upwards in the working dir
			for i := len(dirSplit) - 1; i >= 0; i-- {
				dirName = dirSplit[i]
				entry, exists = conf.Repositories[dirName]
				if exists {
					break
				}
			}
			if !exists {
				fmt.Printf("current directory has not been added to sesame. Do you want to add it? [Y/n]: ")
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				if scanner.Err() != nil {
					_, _ = fmt.Fprintf(os.Stderr, "could not read response: %v\n", err)
				}
				response := scanner.Text()
				switch strings.ToLower(response) {
				case "y", "yes", "":
					repo, err := git.FindRepo(workingDir)
					if err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "current directory is not a repository: %v\n", err)
						os.Exit(1)
					}
					_, err = config.StoreRepositories([]*git.Repository{repo})
					if err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "could not store repository in config: %v\n", err)
						os.Exit(1)
					}
					fmt.Printf("successfully added current directory\n")
					conf, err = config.GetConfig()
					if err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "could not read the config file: %v\n", err)
						os.Exit(1)
					}
					entry = conf.Repositories[repo.Name]
				default:
					os.Exit(0)
				}
			}
			err = browser.OpenURL(entry.Url)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "could not open repository: %v\n", err)
				os.Exit(1)
			}
		} else {
			entry, exists := conf.Repositories[args[0]]
			if !exists {
				_, _ = fmt.Fprintf(os.Stderr, "unknown repository, please add it first: name=%s\n", args[0])
				os.Exit(1)
			}
			err = browser.OpenURL(entry.Url)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "could not open repository: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
