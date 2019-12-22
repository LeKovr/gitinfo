package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"strings"

	"github.com/lekovr/gitinfo"
)

var (
	// Actual version value will be set at build time
	version = "0.0-dev"

	// ErrPathEmpty returned after showing requested help
	ErrPathEmpty = errors.New("required path value is empty")
)

func run(exitFunc func(code int)) {
	var err error
	var cfg *Config
	defer func() { shutdown(exitFunc, err) }()
	cfg, err = setupConfig()
	if err != nil {
		return
	}
	if cfg.Debug {
		log.Printf("gitinfo %s\n", version)
	}
	path := cfg.Args.Path

	if path == "" {
		err = ErrPathEmpty
		return
	}

	if !strings.HasSuffix(path, "/") {
		err = gitinfo.New(cfg.GitInfo).Write(path, nil)
		return
	}

	var (
		d      *os.File
		files  []os.FileInfo
		linked string
		fi     os.FileInfo
	)

	d, err = os.Open(path)
	if err != nil {
		return
	}
	defer d.Close()
	files, err = d.Readdir(-1)
	if err != nil {
		return
	}
	gi := gitinfo.New(cfg.GitInfo)

	for _, file := range files {
		if file.Mode().IsDir() || file.Mode()&os.ModeSymlink != 0 {
			src := filepath.Join(path, file.Name())
			if !file.Mode().IsDir() {

				// resolve symlink
				linked, err = filepath.EvalSymlinks(src)
				if err != nil {
					return
				}
				fi, err = os.Lstat(linked)
				if err != nil {
					return
				}
				if !fi.Mode().IsDir() {
					continue
				}
			}

			err = gi.Write(src, nil)
			if err != nil {
				return
			}
		}
	}
}

// exit after deferred cleanups have run
func shutdown(exitFunc func(code int), e error) {
	if e != nil {
		var code int
		switch e {
		case ErrGotHelp:
			code = 3
		case ErrBadArgs:
			code = 2
		default:
			code = 1
			log.Printf("Run error: %#v", e)
		}
		exitFunc(code)
	}
}
