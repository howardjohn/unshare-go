// Package userns makes the process enter a new network namespace.
package userns

/*
#cgo CFLAGS: -Wall
#define _GNU_SOURCE
#include <stdlib.h>
#include <unistd.h>
#include <sched.h>
#include <stdio.h>
#include <fcntl.h>
#include <string.h>
#include <errno.h>

int originalUid = 0;
int originalGid = 0;

__attribute((constructor(101))) void enter_userns(void) {
		originalUid = getuid();
		originalGid = getgid();
    if (unshare(CLONE_NEWUSER) == -1) {
		fprintf(stderr, "Failed to unshare user namespace: %s\n", strerror(errno));
        exit(1);
    }

    int fd = open("/proc/self/uid_map", O_WRONLY);
    if (fd == -1) {
		fprintf(stderr, "Failed to open /proc/self/uid_map: %s\n", strerror(errno));
        exit(1);
    }

    char uid_map[100];
    snprintf(uid_map, sizeof(uid_map), "0 %d 1\n", originalUid);
    if (write(fd, uid_map, sizeof(uid_map)) == -1) {
        close(fd);
        fprintf(stderr, "Failed to write to /proc/self/uid_map: %s\n", strerror(errno));
        exit(1);
    }

    close(fd);
    return;
}
*/
import "C"

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func OriginalUID() uint32 {
	return uint32(C.originalUid)
}

func OriginalGID() uint32 {
	return uint32(C.originalGid)
}

// WriteUserMap builds a map of Host UID -> Namespace UID
// Example:
//
//	WriteUserMap(map[uint32]uint32{userns.OriginalUID(): 0, 1234: 1234})
func WriteUserMap(mapping map[uint32]uint32) error {
	return writeMap("/proc/self/uid_map", mapping)
}

// WriteGroupMap builds a map of Host GID -> Namespace GID
// Example:
//
//	WriteGroupMap(map[uint32]uint32{userns.OriginalGID(): 0, 1234: 1234})
func WriteGroupMap(mapping map[uint32]uint32) error {
	// write deny in setgroups to disable setgroup(2) and enable writing to gid_map
	err := os.WriteFile("/proc/self/setgroups", []byte("deny\n"), 0o644)
	if err != nil {
		return fmt.Errorf("failed to deny setgroups (%v)", err)
	}

	return writeMap("/proc/self/gid_map", mapping)
}

func writeMap(path string, mapping map[uint32]uint32) error {
	lines := []string{}
	for h, c := range mapping {
		lines = append(lines, fmt.Sprintf("%d %d 1", c, h))
	}
	slices.Sort(lines)
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644)
}
