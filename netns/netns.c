#define _GNU_SOURCE

#include <sched.h>
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <fcntl.h>
#include <string.h>

__attribute((constructor(102))) void enter_netns(void) {
    if (unshare(CLONE_NEWNET) == -1) {
        perror("unshared call failed");
        exit(1);
    }
}
