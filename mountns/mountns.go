// Package mountns makes the process enter a new mount namespace.
package mountns

/*
#cgo CFLAGS: -Wall
#define _GNU_SOURCE

#include <sched.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>

__attribute((constructor(103))) void enter_mountns(void) {
    if (unshare(CLONE_NEWNS) == -1) {
        fprintf(stderr, "Failed to unshare mount namespace: %s\n", strerror(errno));
        exit(1);
    }
}
*/
import "C"

import (
	"syscall"
)

func BindMount(src, dst string) error {
	// Mount source to target using syscall.Mount
	return syscall.Mount(src, dst, "", syscall.MS_BIND, "")
}
