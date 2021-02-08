package gitinfo_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wojas/genericr"

	"github.com/pgmig/gitinfo"
)

const myRepoSuffix = "pgmig/gitinfo.git"

var (
	testHasGit bool
	log        logr.Logger
)

func TestMain(m *testing.M) {
	// look for git binary once
	if _, err := exec.LookPath("git"); err == nil {
		testHasGit = true
	}
	log = genericr.New(func(e genericr.Entry) {
		fmt.Fprintln(os.Stderr, e.String())
	})
	os.Exit(m.Run())
}

func TestVersionError(t *testing.T) {
	gi := gitinfo.New(log, gitinfo.Config{})
	require.NotNil(t, gi.Version(".", nil))
}

func TestRepositoryError(t *testing.T) {
	gi := gitinfo.New(log, gitinfo.Config{})
	require.NotNil(t, gi.Repository(".", nil))
}

func TestModifiedError(t *testing.T) {
	gi := gitinfo.New(log, gitinfo.Config{})
	require.NotNil(t, gi.Modified(".", nil))
}

func TestMakeErrors(t *testing.T) {
	gi := gitinfo.New(log, gitinfo.Config{})
	tests := []struct {
		name string
		path string
		err  error
	}{
		{"Path is empty", "", gitinfo.ErrPathMustNotBeEmpty},
		{"File is not a dir", "gitinfo.go", gitinfo.ErrPathMustBeDir},
	}
	for _, tt := range tests {
		err := gi.Make(tt.path, nil)
		assert.Equal(t, tt.err, err, tt.name)
	}
	gi = gitinfo.New(log, gitinfo.Config{Root: "."})
	err := gi.Make(".notexists", nil)
	if !errors.Is(err, os.ErrNotExist) {
		// We want different error
		require.NoError(t, err)
	}
}

func TestMake(t *testing.T) {
	gi := gitinfo.New(log, gitinfo.Config{})
	require.NoError(t, gi.Make(".", nil))
}

type testFS struct{}

func (fs testFS) Open(name string) (gitinfo.File, error) { return os.Open(name) }

func TestWrite(t *testing.T) {
	tmpName := "tmpfile"
	gi := gitinfo.New(log, gitinfo.Config{File: tmpName, GitBin: "git", Root: "."})
	err := gi.Write(".", nil)
	require.NoError(t, err)

	data, err := gi.Read(testFS{}, ".")
	os.Remove(tmpName)
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(data.Repository, myRepoSuffix))

	data, err = gi.ReadOrMake(testFS{}, ".")
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(data.Repository, myRepoSuffix))
}

func TestMkTime(t *testing.T) {
	var rv *time.Time
	err := gitinfo.MkTime([]byte("xxx"), rv)
	require.NotNil(t, err)
}
