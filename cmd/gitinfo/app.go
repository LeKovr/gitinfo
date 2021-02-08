package main

import (
	"errors"

	"github.com/pgmig/gitinfo"
)

// Config holds all config vars
type Config struct {
	Flags
	Args struct {
		Path []string `description:"path to repository dir(s)"`
	} `positional-args:"yes" required:"yes"`
	GitInfo gitinfo.Config `group:"GitInfo Options" namespace:"gi"`
}

var (
	// Actual version value will be set at build time
	version = "0.0-dev"

	// ErrPathEmpty returned after showing requested help
	ErrPathEmpty = errors.New("required path value is empty")
)

// Run app and exit via given exitFunc
func Run(exitFunc func(code int)) {
	cfg, err := SetupConfig()
	log := SetupLog(err != nil || cfg.Debug)
	defer func() { Shutdown(exitFunc, err, log) }()
	log.V(1).Info("gitinfo. Fetch git repo info.", "v", version)
	if err != nil || cfg.Version {
		return
	}
	if len(cfg.Args.Path) == 0 {
		err = ErrPathEmpty
		return
	}
	gi := gitinfo.New(log, cfg.GitInfo)
	for _, path := range cfg.Args.Path {
		err = gi.Write(path, nil)
		if err == nil {
			continue
		}
		if err == gitinfo.ErrPathMustBeDir {
			// Silently skip plain files
			err = nil
			continue
		}
		return
	}
}
