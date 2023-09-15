package unshare_test

import (
	"fmt"
	"net"

	_ "github.com/howardjohn/unshare-go/netns"
	_ "github.com/howardjohn/unshare-go/userns"
)

func Example_NetNs() {
	l, err := net.Listen("tcp", "127.0.0.1:80")
	if err != nil {
		panic(err.Error())
	}
	go net.Dial("tcp", "127.0.0.1:80")

	if _, err := l.Accept(); err != nil {
		panic(err.Error())
	}
	fmt.Println("Got connection")
	// Output: Got connection
}
