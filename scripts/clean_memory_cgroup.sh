#!/usr/bin/env bash

# let container.go have time to get information of cgroup/memory/$1
sleep(4)
rmdir /sys/fs/cgroup/memory/$1