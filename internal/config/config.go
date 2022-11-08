package config

import (
	"fmt"
	"os"

	"github.com/hugowangler/sesame/internal/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RepositoryEntry struct {
	Url       string   `yaml:"url"`
	Directory string   `yaml:"directory"`
	Links     []string `yaml:"links"`
}

type Config struct {
	Repositories map[string]RepositoryEntry `yaml:"repositories"`
}

func GetConfig() (*Config, error) {
	var config Config
	config.Repositories = make(map[string]RepositoryEntry)
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %v", err)
	}
	return &config, nil
}

func StoreRepositories(repos []*git.Repository) (int, error) {
	var numStored int
	config, err := GetConfig()
	if err != nil {
		return 0, err
	}
	for _, r := range repos {
		config.Repositories[r.Name] = RepositoryEntry{Url: r.Url(), Directory: r.Directory}
		numStored++
	}
	viper.Set("repositories", config.Repositories)
	err = viper.WriteConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			err = viper.WriteConfigAs(home + "/.config/.sesame.yaml")
			if err != nil {
				return 0, fmt.Errorf("failed to create config file: %v", err)
			}
			return numStored, nil
		}
		return 0, fmt.Errorf("failed to update existing config file: %v", err)
	}
	return numStored, nil
}

func StoreLink(repo string, URL string) (int, error) {
	config, err := GetConfig()
	if err != nil {
		return 0, err
	}
	confRepo, exists := config.Repositories[repo]
	if !exists {
		return 0, fmt.Errorf("repository does not exist in the config file: repo=%s", repo)
	}
	confRepo.Links = append(confRepo.Links, URL)
	err = viper.WriteConfig()
	if err != nil {
		return 0, fmt.Errorf("failed to write config file: %v", err)
	}
	return len(confRepo.Links), nil
}
