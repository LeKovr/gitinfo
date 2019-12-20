package gitinfo

import (
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// GitInfo holds git repository metadata
type GitInfo struct {
	Version    string    `json:"version"`
	Repository string    `json:"repository"`
	Modified   time.Time `json:"modified"`
	Path       string    `json:"-"`
	UseGit     bool      `json:"-"`
}

// New returns git info
func New(path string) (rv *GitInfo, err error) {

	now := time.Now()
	rv = &GitInfo{
		Path:       path,
		Repository: "(none)",
		Modified:   now, // TODO: last change of dir content
		Version:    "v0.0.0-" + now.Format("20060102150405"),
	}

	if _, err := exec.LookPath("git"); err == nil {
		rv.UseGit = true
	}
	if rv.UseGit {
		// error means path is not git repo, so skip them
		_ = Version(path, &rv.Version)
		_ = Repository(path, &rv.Repository)
		_ = Modified(path, &rv.Modified)
	}
	return
}

// Version fills rv with package version from git
func Version(path string, rv *string) error {
	out, err := exec.Command("git", "-C", path, "describe", "--tags", "--always").Output()
	if err != nil {
		return err
	}
	*rv = strings.TrimSuffix(string(out), "\n")
	return nil
}

// Repository fills rv with package repo from git
func Repository(path string, rv *string) error {
	out, err := exec.Command("git", "-C", path, "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return err
	}
	*rv = strings.TrimSuffix(string(out), "\n")
	return nil
}

// Modified fills rv with package last commit timestamp
func Modified(path string, rv *time.Time) error {
	out, err := exec.Command("git", "-C", path, "show", "-s", "--format=format:%ct", "HEAD").Output()
	if err != nil {
		return err
	}
	return MkTime(out, rv)
}

// MkTime converts []byte to time.Time
func MkTime(in []byte, rv *time.Time) error {
	tm, err := strconv.ParseInt(string(in), 10, 64)
	if err != nil {
		return err
	}
	*rv = time.Unix(tm, 0)
	return nil
}
