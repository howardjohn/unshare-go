package unshare_test

import (
	"fmt"
	_ "github.com/howardjohn/unshare-go/userns"
	"os"
)

func Example_UserNs() {
	fmt.Println("Running as user", os.Getuid())
	// Output: Running as user 0
}
