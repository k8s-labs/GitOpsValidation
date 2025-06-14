package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
	gitConfigPath := filepath.Join(dir, ".git", "config")
	data, err := ioutil.ReadFile(gitConfigPath)
	if err != nil {
		return err
	}
	if !containsRemoteURL(string(data), expectedURL) {
		return fmt.Errorf("repo at %s does not match expected remote %s", dir, expectedURL)
	}
	return nil
}

// containsRemoteURL checks if the .git/config contains the expected remote URL
func containsRemoteURL(config, url string) bool {
	return (len(config) > 0 && url != "" && (string(config) == url || (len(config) > 0 && (filepath.Base(url) == filepath.Base(config)))))
}

// CheckoutBranch checks out the specified branch in the given repo directory
func CheckoutBranch(dir, branch string) error {
	cmd := exec.Command("git", "-C", dir, "checkout", branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git checkout failed: %v, output: %s", err, string(output))
	}
	return nil
}

// PullLatest pulls the latest changes from the current branch in the given repo directory
func PullLatest(dir string) error {
	cmd := exec.Command("git", "-C", dir, "pull")
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
