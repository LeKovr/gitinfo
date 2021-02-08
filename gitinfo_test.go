package gitinfo_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wojas/genericr"

	"github.com/pgmig/gitinfo"
)

const myRepoSuffix = "pgmig/gitinfo.git"

var (
	testHasGit bool
	tempDir    string
)

func TestMain(m *testing.M) {
	// look for git binary once
	if _, err := exec.LookPath("git"); err == nil {
		testHasGit = true
	}
	if testHasGit {
		var err error
		tempDir, err = ioutil.TempDir("/tmp/", "gitinfo-test")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(tempDir)
	}
	os.Exit(m.Run())
}

func TestVersion(t *testing.T) {
	if !testHasGit {
		t.Skip("git not found, skipping")
	}
	var rv string
	if err := gitinfo.Version(".", &rv); err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.NotEqual(t, rv, "")
	if err := gitinfo.Version(tempDir, &rv); err == nil {
		t.Fatalf("Call must return error")
	}
}

func TestRepository(t *testing.T) {
	if !testHasGit {
		t.Skip("git not found, skipping")
	}
	var rv string
	if err := gitinfo.Repository(".", &rv); err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.NotEqual(t, rv, "")
	if err := gitinfo.Repository(tempDir, &rv); err == nil {
		t.Fatalf("Call must return error")
	}
}

func TestModified(t *testing.T) {
	if !testHasGit {
		t.Skip("git not found, skipping")
	}
	var rv, zeroTm time.Time
	if err := gitinfo.Modified(".", &rv); err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.NotEqual(t, rv, zeroTm)
	if err := gitinfo.Modified(tempDir, &rv); err == nil {
		t.Fatalf("Call must return error")
	}
}

func TestMakeErrors(t *testing.T) {
	log := genericr.New(func(e genericr.Entry) {
		fmt.Fprintln(os.Stderr, e.String())
	})
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
		// Not error we want
		require.NoError(t, err)
	}

}

type testFS struct{}

func (fs testFS) Open(name string) (gitinfo.File, error) { return os.Open(name) }

func TestWrite(t *testing.T) {
	log := genericr.New(func(e genericr.Entry) {
		fmt.Fprintln(os.Stderr, e.String())
	})
	tmpName := "tmpfile"
	gi := gitinfo.New(log, gitinfo.Config{File: tmpName})
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
