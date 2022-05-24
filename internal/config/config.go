package config

import (
	"fmt"
	"github.com/hugowangler/sesame/internal/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Repositories map[string]string `yaml:"repositories"`
}

func StoreRepositories(repos []*git.Repository) (int, error) {
	var numStored int
	var config Config
	config.Repositories = make(map[string]string)
	err := viper.Unmarshal(&config)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal config file: err=%v", err)
	}
	for _, r := range repos {
		if _, exists := config.Repositories[r.Name]; !exists {
			config.Repositories[r.Name] = r.Url()
			numStored++
		}
	}
	viper.Set("repositories", config.Repositories)
	err = viper.WriteConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			err = viper.WriteConfigAs(home + "/.config/.sesame.yaml")
			if err != nil {
				return 0, fmt.Errorf("failed to create config file: err=%v\n", err)
			}
			return numStored, nil
		}
		return 0, fmt.Errorf("failed to update existing config file: err=%v\n", err)
	}
	return numStored, nil
}
