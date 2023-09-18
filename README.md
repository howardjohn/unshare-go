# `unshare-go`

[![Go Reference](https://pkg.go.dev/badge/github.com/howardjohn/unshare-go.svg)](https://pkg.go.dev/github.com/howardjohn/unshare-go)

`unshare-go` is a small module to make use of the [`unshare`](https://man7.org/linux/man-pages/man2/unshare.2.html) syscall.

This allows creating new namespaces for a process, which can be used for things like running networking actions that typically require high privileges, etc.

## How it works

`unshare` has interactions with threads. Executions either impact on the current thread, or cannot run at all in a threaded program.
This is problematic as all Go programs are threaded.

To get around this, *CGO is used* to execute the `unshare` actions extremely early in the process initialization, before Go's runtime has spawned any threads.
As a result, usage of this library is not done by calling functions, but rather importing the packages.
In many ways this is like `init()` functions, but called even earlier.

As a result of this design, the library is not very flexible.

Currently, only user and network namespaces are supported, although others could be added.

## Usage

To use the library, simply import the package:

```go
import (
	// Create a new user namespace. This will map the current UID to 0.
	_ "github.com/howardjohn/unshare-go/userns"
	// Create a new network namespace. This will have the 'lo' interface ready but nothing else.
	_ "github.com/howardjohn/unshare-go/netns"
	// Create a new mount namespace.
	_ "github.com/howardjohn/unshare-go/mountns"
)
```

Namespace creation is ordered:

1. User namespace
2. Network namespace
3. Mount namespace

That is, the network namespace is created inside the user namespace.

An end to end example:

```go
package main

import (
	"log"
	"net"
	"os"
	
	_ "github.com/howardjohn/unshare-go/userns"
	_ "github.com/howardjohn/unshare-go/netns"
)

func main() {
	log.Println(os.Getuid()) // Returns 0, always
	net.Listen("tcp", "127.0.0.1:80") // This works, even when the process is not run as root.
}
```
