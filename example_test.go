package gitinfo

import (
	"fmt"
	"strings"
	"time"
)

func ExampleNew() {

	data, err := New("cmd/")
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	fmt.Printf("%v\n%v\n%v\n",
		data.Modified != time.Time{},
		data.Version != "",
		strings.HasSuffix(data.Repository, "LeKovr/gitinfo.git"),
	)
	// Output:
	// true
	// true
	// true
}
