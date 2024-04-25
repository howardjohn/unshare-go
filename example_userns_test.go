package unshare_test

import (
	"fmt"
	"os"

	_ "github.com/howardjohn/unshare-go/userns"
)

func Example_UserNs() {
	fmt.Println("Running as user", os.Getuid())
	fmt.Println("Running as group", os.Getgid())
	// Output:
	// Running as user 0
	// Running as group 0
}
