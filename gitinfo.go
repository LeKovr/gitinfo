// Package gitinfo used for generate and read gitinfo.json file
// which contains git repository data like this:
// ```
// {
//   "version": "v0.33-1-g4f4575a",
//   "repository": "https://github.com/pgmig-sql/pgmig.git",
//   "modified": "2019-12-24T02:44:51+03:00"
// }
// ```
// This file is intended to be included in embedded FS
package gitinfo

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Config holds all config vars
type Config struct {
	Debug bool   `long:"debug" description:"Show debug data"`
	File  string `long:"file" default:"gitinfo.json" description:"GitInfo json filename"`
	// Root may be hardcoded if app uses embedded FS, so do not show it in help
	Root string
}

// GitInfo holds git repository metadata
type GitInfo struct {
	Version    string    `json:"version"`
	Repository string    `json:"repository"`
	Modified   time.Time `json:"modified"`
}

// Service holds service data
type Service struct {
	Config
	GI *GitInfo
}

// New returns service object with config
func New(cfg Config) *Service {
	return &Service{Config: cfg}
}

// Make prepares GitInfo data
func (srv Service) Make(path string, gi *GitInfo) error {
	if srv.Config.Root != "" {
		path = filepath.Join(srv.Config.Root, path)
	}
	// check for dir exists
	info, err := os.Stat(path)
	if err != nil { //os.IsNotExist(err) {
		return errors.Wrap(err, "Path is not available")
	}
	if !info.IsDir() {
		return errors.Wrap(err, "Path is not directory")
	}

	useGit := false

	// check for git bin available
	if _, err := exec.LookPath("git"); err == nil {
		// check for dir is a repo
		if err = Repository(path, &gi.Repository); err == nil {
			useGit = true
		}
	}

	now := time.Now()

	if useGit {
		if err = Version(path, &gi.Version); err != nil {
			// Repo has no tags, generate own
			// TODO: log.warn
			gi.Version = "v0.0.0-" + now.Format("20060102150405")
		}
		if err = Modified(path, &gi.Modified); err != nil {
			// TODO: log.warn
			gi.Modified = now // TODO: last change of dir content
		}
	} else {
		abs, err := filepath.Abs(path)
		if err != nil {
			return errors.Wrap(err, "Resolve relative path")
		}
		gi.Repository = "file://" + abs
		gi.Version = "v0.0.0-" + now.Format("20060102150405")
		gi.Modified = now // TODO: last change of dir content
	}
	return nil
}

// Write saves GitInfo data to file, prepare it if none given
func (srv Service) Write(path string, gi *GitInfo) error {

	if gi == nil {
		if srv.Config.Debug {
			log.Printf("Looking in %s", path)
		}
		gi = &GitInfo{}
		err := srv.Make(path, gi)
		if err != nil {
			return err
		}
	}
	fn := filepath.Join(path, srv.Config.File)
	f, err := os.Create(fn)
	if err != nil {
		return errors.Wrap(err, "Create gitinfo file")
	}
	defer f.Close()

	out, err := json.MarshalIndent(gi, "", "  ")
	if err != nil {
		return errors.Wrap(err, "Create gitinfo json")
	}

	_, err = f.WriteString(string(out) + "\n") //ioutil.WriteFile(p, out, os.FileMode(mode))
	if err != nil {
		return errors.Wrap(err, "Write gitinfo file")
	}
	return nil
}

// Version fills rv with package version from git
func Version(path string, rv *string) error {
	out, err := exec.Command("git", "-C", path, "describe", "--tags", "--always").Output()
	if err != nil {
		return errors.Wrap(err, "Git describe")
	}
	*rv = strings.TrimSuffix(string(out), "\n")
	return nil
}

// Repository fills rv with package repo from git
func Repository(path string, rv *string) error {
	out, err := exec.Command("git", "-C", path, "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return errors.Wrap(err, "Git repo")
	}
	*rv = strings.TrimSuffix(string(out), "\n")
	return nil
}

// Modified fills rv with package last commit timestamp
func Modified(path string, rv *time.Time) error {
	out, err := exec.Command("git", "-C", path, "show", "-s", "--format=format:%ct", "HEAD").Output()
	if err != nil {
		return errors.Wrap(err, "Git show")
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

// File is an interface for FileSystem.Open func
type File interface {
	io.Closer
	io.Reader
	io.Seeker
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}

// FileSystem holds all of used filesystem access methods
type FileSystem interface {
	Walk(root string, walkFn filepath.WalkFunc) error
	Open(name string) (File, error)
}

// Read reads GitInfo data from file
func (srv Service) Read(fs FileSystem, path string) (*GitInfo, error) {

	fn := filepath.Join(path, srv.Config.File)
	file, err := fs.Open(fn)
	if err != nil {
		return nil, errors.Wrap(err, "Open gitinfo file")
	}
	defer file.Close()

	js, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "Read gitinfo file")
	}

	gi := GitInfo{}
	err = json.Unmarshal(js, &gi)
	if err != nil {
		return nil, errors.Wrap(err, "Parse gitinfo file")
	}
	return &gi, nil
}

// ReadOrMake reads GitInfo data from file or makes it from git
func (srv Service) ReadOrMake(fs FileSystem, path string) (*GitInfo, error) {

	gi, err := srv.Read(fs, path)
	if err != nil {
		gi = &GitInfo{}
		err = srv.Make(path, gi)
	}
	return gi, err
}
