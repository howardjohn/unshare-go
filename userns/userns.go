// Package userns makes the process enter a new network namespace.
package userns

/*
#cgo CFLAGS: -Wall
#define _GNU_SOURCE
#include <stdlib.h>
#include <unistd.h>
#include <sched.h>

 int originalUid = 0;
 int originalGid = 0;

__attribute((constructor(101))) void enter_userns(void) {
	originalUid = getuid();
	originalGid = getgid();
    if (unshare(CLONE_NEWUSER) == -1) {
        exit(1);
    }
}

*/
import "C"

import (
	"fmt"
	"os"
	"strings"

	"slices"
)

func init() {
	err := WriteMap("/proc/self/uid_map", map[uint32]uint32{
		OriginalUID(): 0,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create uid mapping: %v", err).Error())
	}

	// write deny in setgroups to disable setgroup(2) and enable writing to gid_map
	data := []byte("deny\n")
	err = os.WriteFile("/proc/self/setgroups", data, 0644)
	if err != nil {
		fmt.Println("Error writing to setgroups:", err)
	}

	err = WriteMap("/proc/self/gid_map", map[uint32]uint32{
		OriginalGID(): 0,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create gid mapping: %v", err).Error())
	}
}

func OriginalUID() uint32 {
	return uint32(C.originalUid)
}

func OriginalGID() uint32 {
	return uint32(C.originalGid)
}

// WriteMap builds a map of Host UID/GID -> Namespace UID/GID
// Example:
//
//	WriteMap(map[uint32]uint32{userns.OriginalUID: 0, 1234: 1234})
func WriteMap(path string, mapping map[uint32]uint32) error {
	lines := []string{}
	for h, c := range mapping {
		lines = append(lines, fmt.Sprintf("%d %d 1", c, h))
	}
	slices.Sort(lines)
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644)
}
