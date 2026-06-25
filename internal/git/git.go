package git

import (
	"bytes"
	"context"
	"fmt"
	"gopkg.in/ini.v1"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// see https://www.debuggex.com/r/fFggA8Uc4YYKjl34
const repoRegex = `((git@|https://|http://)(?P<host>[\w\.@]+)(/|:))(?P<owner>[\w,\-,\_]+)/(?P<name>[\w,\-,\_]+)(.git){0,1}((/){0,1})`

var subexps = []string{
	"host",
	"owner",
	"name",
}

type Repository struct {
	host      string
	owner     string
	Name      string
	Directory string
}

func (r *Repository) Url() string {
	return fmt.Sprintf("https://%s/%s/%s", r.host, r.owner, r.Name)
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
		repo := parseGitConfig(filepath.Join(p, "config"), filepath.Dir(p))
		if repo != nil {
			repos = append(repos, repo)
		}
	}

	return repos, nil
}

func FindRepo(path string) (*Repository, error) {
	topLevel, err := runGit(path, "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, fmt.Errorf("could not find git repository: path=%s, err=%v", path, err)
	}

	url, err := runGit(path, "config", "--local", "--get", "remote.origin.url")
	if err != nil {
		return nil, fmt.Errorf("could not find remote origin URL: path=%s, err=%v", path, err)
	}

	repo := parseRepositoryURL(url, topLevel, path)
	if repo == nil {
		return nil, fmt.Errorf("could not parse remote origin URL: path=%s, url=%s", path, url)
	}

	return repo, nil
}

func runGit(path string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", path}, args...)...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("git command timed out")
		}
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		return "", fmt.Errorf("%s", message)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func parseGitConfig(configPath string, repoPath string) *Repository {
	cfg, err := ini.Load(configPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not read .git/config file: path=%s, err=%v\n", configPath, err)
		return nil
	}
	url := cfg.Section("remote \"origin\"").Key("url").String()
	if url == "" {
		_, _ = fmt.Fprintf(os.Stderr, "could not find a URL in .git/config, ignoring repository: path=%s\n", configPath)
		return nil
	}

	return parseRepositoryURL(url, repoPath, configPath)
}

func parseRepositoryURL(url string, repoPath string, sourcePath string) *Repository {
	re := regexp.MustCompile(repoRegex)
	submatches := re.FindStringSubmatch(url)
	m := make(map[string]string, len(subexps))
	for _, s := range subexps {
		if len(submatches) < len(subexps) {
			_, _ = fmt.Fprintf(os.Stderr, "could not parse URL, ignoring repository: path=%s, url=%s\n", sourcePath, url)
			return nil
		}
		match := submatches[re.SubexpIndex(s)]
		if match == "" {
			_, _ = fmt.Fprintf(
				os.Stderr,
				"could not regex match required repository information from the found URL: path=%s url=%s\n",
				sourcePath,
				url,
			)
			return nil
		}
		m[s] = strings.ToLower(match)
	}

	return &Repository{
		host:      m[subexps[0]],
		owner:     m[subexps[1]],
		Name:      strings.ToLower(m[subexps[2]]),
		Directory: filepath.Base(filepath.Clean(repoPath)),
	}
}
