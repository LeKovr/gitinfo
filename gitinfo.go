package gitinfo

import (
	"os/exec"
	"strconv"
	"strings"
	"time"
)

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
	tm, err := strconv.ParseInt(string(out), 10, 64)
	if err != nil {
		return err
	}
	*rv = time.Unix(tm, 0)
	return nil
}