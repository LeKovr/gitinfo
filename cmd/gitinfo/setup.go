package main

import (
	"errors"
	//	"fmt"

	"github.com/pgmig/gitinfo"
	"github.com/jessevdk/go-flags"
)

// Config holds all config vars
type Config struct {
	Debug bool `long:"debug" description:"Show debug data"`
	Args  struct {
		Path string `description:"path to repository dir(s)"`
	} `positional-args:"yes" required:"yes"`
	GitInfo gitinfo.Config `group:"GitInfo Options" namespace:"gi"`
}

var (
	// ErrGotHelp returned after showing requested help
	ErrGotHelp = errors.New("help printed")
	// ErrBadArgs returned after showing command args error message
	ErrBadArgs = errors.New("option error printed")
)

// setupConfig loads Config fields
func setupConfig(args ...string) (*Config, error) {
	cfg := &Config{}
	p := flags.NewParser(cfg, flags.Default) //  HelpFlag | PrintErrors | PassDoubleDash
	var err error
	if len(args) == 0 {
		_, err = p.Parse()
	} else {
		_, err = p.ParseArgs(args)
	}
	if err != nil {
		//fmt.Printf("Args error: %#v", err)
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			return nil, ErrGotHelp
		}
		return nil, ErrBadArgs
	}
	return cfg, nil
}
