package git

import (
	"fmt"
	"gopkg.in/ini.v1"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// see https://www.debuggex.com/r/fFggA8Uc4YYKjl34
const repoRegex = `((git@|https://)(?P<host>[\w\.@]+)(/|:))(?P<owner>[\w,\-,\_]+)/(?P<name>[\w,\-,\_]+)(.git){0,1}((/){0,1})`

var subexps = []string{
	"host",
	"owner",
	"name",
}

type Repository struct {
	Host  string
	Owner string
	Name  string
}

func (r *Repository) Url() string {
	return fmt.Sprintf("https://%s/%s/%s", r.Host, r.Owner, r.Name)
}

func FindRepos(path string) ([]*Repository, error) {
	var gitConfigPaths []string
	err := filepath.WalkDir(
		path, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() && d.Name() == ".git" {
				gitConfigPaths = append(gitConfigPaths, path)
				return fs.SkipDir
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	var repos []*Repository
	for _, p := range gitConfigPaths {
		repo := parseGitConfig(p)
		if repo != nil {
			repos = append(repos, repo)
		}
	}

	return repos, nil
}

func parseGitConfig(path string) *Repository {
	cfg, err := ini.Load(path + "/config")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not read .git/config file: path=%s, err=%v\n", path, err)
		return nil
	}
	url := cfg.Section("remote \"origin\"").Key("url").String()
	if url == "" {
		_, _ = fmt.Fprintf(os.Stderr, "could not find a URL in .git/config, ignoring repository: path=%s\n", path)
		return nil
	}

	re := regexp.MustCompile(repoRegex)
	submatches := re.FindStringSubmatch(url)
	m := make(map[string]string, len(subexps))
	for _, s := range subexps {
		match := submatches[re.SubexpIndex(s)]
		if match == "" {
			_, _ = fmt.Fprintf(
				os.Stderr,
				"could not regex match required repository information from the found URL: path=%s url=%s\n",
				path,
				url,
			)
			return nil
		}
		m[s] = strings.ToLower(match)
	}
	return &Repository{
		Host:  m[subexps[0]],
		Owner: m[subexps[1]],
		Name:  m[subexps[2]],
	}
}
