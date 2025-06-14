package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CloneOptions struct {
	RepoURL string
	User    string
	PAT     string
	Branch  string
	Dir     string // Target directory
}

// CloneRepo clones the GitOps repository using provided credentials
func CloneRepo(opts CloneOptions) error {
	url := opts.RepoURL

	if opts.User != "" && opts.PAT != "" {
		// Insert credentials into URL if needed
		url = fmt.Sprintf("https://%s:%s@%s", opts.User, opts.PAT, url[8:])
	}

	args := []string{"clone", "--branch", opts.Branch, url, opts.Dir}
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %v, output: %s", err, string(output))
	}
	return nil
}

// VerifyRepo checks if the existing repo directory matches the configured repo URL
func VerifyRepo(dir, expectedURL string) error {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git remote get-url failed: %v, output: %s", err, string(output))
	}

	// Compare the actual URL with expected URL (trimming whitespace)
	actualURL := string(output)
	actualURL = strings.TrimSpace(actualURL)

	// Compare URLs case-insensitively
	if !strings.EqualFold(actualURL, expectedURL) {
		return fmt.Errorf("repo at %s has remote %s, expected %s", dir, actualURL, expectedURL)
	}

	return nil
}

// CheckoutBranch checks out the specified branch in the given repo directory
func CheckoutBranch(branch string) error {
	cmd := exec.Command("git", "checkout", branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git checkout failed: %v, output: %s", err, string(output))
	}
	return nil
}

// PullLatest pulls the latest changes from the current branch in the given repo directory
func PullLatest(dir string) error {
	cmd := exec.Command("git", "pull")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed: %v, output: %s", err, string(output))
	}
	return nil
}

// ChangeToPath changes the working directory to the specified path within the repo
func ChangeToPath(repoDir, subPath string) error {
	target := filepath.Join(repoDir, subPath)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return fmt.Errorf("specified path does not exist: %s", target)
	}
	return os.Chdir(target)
}

// ExtractRepoName extracts the repository name from a Git URL
func ExtractRepoName(url string) string {
	// Handle both HTTPS and SSH urls
	base := filepath.Base(url)

	// Remove .git suffix if present
	if len(base) > 4 && base[len(base)-4:] == ".git" {
		base = base[:len(base)-4]
	}

	return base
}
