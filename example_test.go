package gitinfo

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"strings"
	"time"
)

func ExampleNew() {

	cfg := Config{}
	// Fill config with default values
	p := flags.NewParser(&cfg, flags.Default|flags.IgnoreUnknown)
	_, err := p.Parse()
	//        require.NoError(ss.T(), err)

	var gi GitInfo
	err = New(cfg).Make("cmd/", &gi)
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
