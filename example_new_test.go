package gitinfo_test

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/wojas/genericr"

	"github.com/pgmig/gitinfo"
)

func ExampleNew() {
	log := genericr.New(func(e genericr.Entry) {
		fmt.Fprintln(os.Stderr, e.String())
	})
	cfg := gitinfo.Config{}
	// Fill config with default values
	p := flags.NewParser(&cfg, flags.Default|flags.IgnoreUnknown)
	_, err := p.Parse()
	if err != nil {
		log.Error(err, "Config")
		os.Exit(1)
	}

	var gi gitinfo.GitInfo
	err = gitinfo.New(log, cfg).Make("cmd/", &gi)
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	fmt.Printf("%v\n%v\n%v\n",
		gi.Modified != time.Time{},
		gi.Version != "",
		strings.HasSuffix(gi.Repository, "pgmig/gitinfo.git"),
	)
	// Output:
	// true
	// true
	// true
}
