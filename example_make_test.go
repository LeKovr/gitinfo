package gitinfo_test

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wojas/genericr"

	"github.com/pgmig/gitinfo"
)

func ExampleService_Make() {
	log := genericr.New(func(e genericr.Entry) {
		fmt.Fprintln(os.Stderr, e.String())
	})
	gi := gitinfo.GitInfo{}
	err := gitinfo.New(log, gitinfo.Config{GitBin: "git"}).Make("cmd/", &gi)
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
