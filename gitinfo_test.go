package gitinfo

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testHasGit bool
var tempDir string

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags

	if _, err := exec.LookPath("git"); err == nil {
		testHasGit = true
	}

	var err error

	if testHasGit {
		tempDir, err = ioutil.TempDir("/tmp/", "gitinfo-test")
		if err != nil {
			log.Fatal(err)
		}
	}

	defer os.Exit(m.Run())
	if testHasGit {
		defer os.RemoveAll(tempDir)
	}
}

func TestVersion(t *testing.T) {
	if !testHasGit {
		t.Skip("git not found, skipping")
	}

	var rv string
	if err := Version(".", &rv); err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.NotEqual(t, rv, "")

	if err := Version(tempDir, &rv); err == nil {
		t.Fatalf("Call must return error")
	}
}

func TestRepository(t *testing.T) {
	if !testHasGit {
		t.Skip("git not found, skipping")
	}

	var rv string
	if err := Repository(".", &rv); err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.NotEqual(t, rv, "")

	if err := Repository(tempDir, &rv); err == nil {
		t.Fatalf("Call must return error")
	}
}

func TestModified(t *testing.T) {
	if !testHasGit {
		t.Skip("git not found, skipping")
	}

	var rv, zeroTm time.Time
	if err := Modified(".", &rv); err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.NotEqual(t, rv, zeroTm)

	if err := Modified(tempDir, &rv); err == nil {
		t.Fatalf("Call must return error")
	}
}

func TestMkTime(t *testing.T) {
	var rv *time.Time
	err := MkTime([]byte("xxx"), rv)

	assert.NotNil(t, err)
}
