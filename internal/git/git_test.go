package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindRepoFromGitDirectory(t *testing.T) {
	repoPath := newRepo(t, "Sesame")

	repo, err := FindRepo(repoPath)
	if err != nil {
		t.Fatalf("FindRepo returned error: %v", err)
	}

	assertRepo(t, repo, "sesame", "Sesame", "https://github.com/hugowangler/sesame")
}

func TestFindRepoFromNestedDirectory(t *testing.T) {
	repoPath := newRepo(t, "Sesame")
	nestedPath := filepath.Join(repoPath, "cmd", "nested")
	if err := os.MkdirAll(nestedPath, 0755); err != nil {
		t.Fatalf("could not create nested directory: %v", err)
	}

	repo, err := FindRepo(nestedPath)
	if err != nil {
		t.Fatalf("FindRepo returned error: %v", err)
	}

	assertRepo(t, repo, "sesame", "Sesame", "https://github.com/hugowangler/sesame")
}

func TestFindRepoFromLinkedWorktree(t *testing.T) {
	_, worktreePath := newRepoWithWorktree(t)

	repo, err := FindRepo(worktreePath)
	if err != nil {
		t.Fatalf("FindRepo returned error: %v", err)
	}

	assertRepo(t, repo, "sesame", "test1", "https://github.com/hugowangler/sesame")
}

func TestFindRepoFromLinkedWorktreeNestedDirectory(t *testing.T) {
	_, worktreePath := newRepoWithWorktree(t)
	nestedPath := filepath.Join(worktreePath, "internal", "git")
	if err := os.MkdirAll(nestedPath, 0755); err != nil {
		t.Fatalf("could not create nested directory: %v", err)
	}

	repo, err := FindRepo(nestedPath)
	if err != nil {
		t.Fatalf("FindRepo returned error: %v", err)
	}

	assertRepo(t, repo, "sesame", "test1", "https://github.com/hugowangler/sesame")
}

func TestFindRepoFromNonRepository(t *testing.T) {
	if _, err := FindRepo(t.TempDir()); err == nil {
		t.Fatalf("FindRepo returned nil error for non-repository path")
	}
}

func newRepoWithWorktree(t *testing.T) (string, string) {
	t.Helper()
	tmp := t.TempDir()
	repoPath := initRepo(t, filepath.Join(tmp, "sesame"))
	worktreePath := filepath.Join(tmp, "sesame__worktrees", "test1")
	runTestGit(t, repoPath, "worktree", "add", "-b", "test1", worktreePath, "HEAD")
	return repoPath, worktreePath
}

func newRepo(t *testing.T, dirName string) string {
	t.Helper()
	return initRepo(t, filepath.Join(t.TempDir(), dirName))
}

func initRepo(t *testing.T, repoPath string) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git executable not found")
	}
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("could not create repo directory: %v", err)
	}

	runTestGit(t, repoPath, "init")
	runTestGit(t, repoPath, "config", "user.email", "test@example.com")
	runTestGit(t, repoPath, "config", "user.name", "Test User")
	runTestGit(t, repoPath, "remote", "add", "origin", "git@github.com:HugoWangler/Sesame.git")
	writeFile(t, filepath.Join(repoPath, "README.md"), "# Sesame\n")
	runTestGit(t, repoPath, "add", "README.md")
	runTestGit(t, repoPath, "commit", "-m", "initial commit")

	return repoPath
}

func runTestGit(t *testing.T, path string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", path}, args...)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
}

func writeFile(t *testing.T, path string, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("could not create parent directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatalf("could not write file: %v", err)
	}
}

func assertRepo(t *testing.T, repo *Repository, name string, directory string, url string) {
	t.Helper()
	if repo == nil {
		t.Fatalf("repo is nil")
	}
	if repo.Name != name {
		t.Fatalf("repo.Name = %q, want %q", repo.Name, name)
	}
	if repo.Directory != directory {
		t.Fatalf("repo.Directory = %q, want %q", repo.Directory, directory)
	}
	if repo.Url() != url {
		t.Fatalf("repo.Url() = %q, want %q", repo.Url(), url)
	}
}
