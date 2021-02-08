// Package gitinfo used for generate and read gitinfo.json file
// which contains git repository data like this:
//  {
//    "version": "v0.33-1-g4f4575a",
//    "repository": "https://github.com/pgmig-sql/pgmig.git",
//    "modified": "2019-12-24T02:44:51+03:00"
//  }
// This file is intended to be included in embedded FS
package gitinfo

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

// Config holds all config vars
type Config struct {
	Debug  bool   `long:"debug" description:"Show debug data"`
	File   string `long:"file" default:"gitinfo.json" description:"GitInfo json filename"`
	GitBin string `long:"gitbin" default:"git" description:"Git binary name"`
	// Root may be hardcoded if app uses embedded FS, so do not showed in help
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
	Config Config
	Log    logr.Logger
	GI     *GitInfo
	useGit bool
}

var (
	// ErrPathMustNotBeEmpty raised if given path is empty
	ErrPathMustNotBeEmpty = errors.New("Path must not be empty")
	// ErrPathMustBeDir raised if given path is not directory
	ErrPathMustBeDir = errors.New("Path must be a directory")
)

// New returns service object with config
func New(log logr.Logger, cfg Config) *Service {
	var useGit bool
	if _, err := exec.LookPath(cfg.GitBin); err != nil {
		log.Error(err, "No git binary found")
	} else {
		useGit = true
	}
	return &Service{Config: cfg, Log: log, useGit: useGit}
}

// Make prepares GitInfo data
func (srv Service) Make(path string, gi *GitInfo) error {
	if srv.Config.Root != "" {
		path = filepath.Join(srv.Config.Root, path)
	}
	if path == "" {
		return ErrPathMustNotBeEmpty
	}
	isDir, err := fileIsDirOrDirLink(path)
	if err != nil {
		return errors.Wrap(err, "Path is not available")
	}
	if !isDir {
		return ErrPathMustBeDir
	}
	useGit := false
	if gi == nil {
		gi = &GitInfo{}
	}
	log := srv.Log.WithValues("path", path)
	// check for git bin available
	if srv.useGit {
		// check for dir is a repo
		if err = srv.Repository(path, &gi.Repository); err == nil {
			useGit = true
		}
	}
	now := time.Now() // MAYBE: last change of dir content?
	if !useGit {
		log.Info("git is not available")
		abs, err := filepath.Abs(path)
		if err != nil {
			return errors.Wrap(err, "Resolve relative path")
		}
		gi.Repository = "file://" + abs
		gi.Version = "v0.0.0-" + now.Format("20060102150405")
		gi.Modified = now
		return nil
	}

	if err = srv.Version(path, &gi.Version); err != nil {
		// Repo has no tags, generate own
		log.Info("repo tag not found")
		gi.Version = "v0.0.0-" + now.Format("20060102150405")
	}
	if err = srv.Modified(path, &gi.Modified); err != nil {
		log.Info("set Modified = now")
		gi.Modified = now
	}
	return nil
}

// Write saves GitInfo data to file, prepare it if none given
func (srv Service) Write(path string, gi *GitInfo) error {
	srv.Log.V(1).Info("Write gitinfo", "file", path)
	if gi == nil {
		srv.Log.V(1).Info("Fetching git metadata")
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
func (srv Service) Version(path string, rv *string) error {
	if srv.Config.Root != "" {
		path = filepath.Join(srv.Config.Root, path)
	}
	out, err := exec.Command(srv.Config.GitBin, "-C", path, "describe", "--tags", "--always").Output()
	if err != nil {
		return errors.Wrap(err, "Git describe")
	}
	*rv = strings.TrimSuffix(string(out), "\n")
	return nil
}

// Repository fills rv with package repo from git
func (srv Service) Repository(path string, rv *string) error {
	if srv.Config.Root != "" {
		path = filepath.Join(srv.Config.Root, path)
	}
	out, err := exec.Command(srv.Config.GitBin, "-C", path, "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return errors.Wrap(err, "Git repo")
	}
	*rv = strings.TrimSuffix(string(out), "\n")
	return nil
}

// Modified fills rv with package last commit timestamp
func (srv Service) Modified(path string, rv *time.Time) error {
	if srv.Config.Root != "" {
		path = filepath.Join(srv.Config.Root, path)
	}
	out, err := exec.Command(srv.Config.GitBin, "-C", path, "show", "-s", "--format=format:%ct", "HEAD").Output()
	if err != nil {
		return errors.Wrap(err, "Git show")
	}
	return MkTime(out, rv)
}

// File is an interface for FileSystem.Open func
type File interface {
	io.Closer
	io.Reader
}

// FileSystem holds all of used filesystem access methods
type FileSystem interface {
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

// ReadOrMake reads GitInfo data from file or makes it from git on the fly
func (srv Service) ReadOrMake(fs FileSystem, path string) (*GitInfo, error) {
	gi, err := srv.Read(fs, path)
	if err != nil {
		gi = &GitInfo{}
		err = srv.Make(path, gi)
	}
	return gi, err
}

// fileIsDirOrDirLink returns true if path is a dir or symlink to dir
func fileIsDirOrDirLink(path string) (bool, error) {
	file, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if file.IsDir() {
		return true, nil
	}
	if file.Mode()&os.ModeSymlink == 0 {
		return false, nil
	}
	// check symlink
	var linkSrc string
	linkDst := filepath.Join(path, file.Name())
	linkSrc, err = filepath.EvalSymlinks(linkDst)
	if err == nil {
		var fi os.FileInfo
		fi, err = os.Lstat(linkSrc)
		if err == nil {
			return fi.IsDir(), nil
		}
	}
	return false, err
}

// MkTime converts to time.Time result of `git show -s --format=format:%ct HEAD`
func MkTime(in []byte, rv *time.Time) error {
	tm, err := strconv.ParseInt(string(in), 10, 64)
	if err != nil {
		return err
	}
	*rv = time.Unix(tm, 0)
	return nil
}
