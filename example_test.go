package gitinfo

import (
	"fmt"
	"time"
)

func ExampleNew() {

	data, err := New("cmd/")
	if err != nil {
		fmt.Printf("%#v\n", err)
	}
	fmt.Printf("%v\n%v\n%s\n", data.Modified != time.Time{}, data.Version != "", data.Repository)
	// Output:
	// true
	// true
	// git@github.com:LeKovr/gitinfo.git
}
