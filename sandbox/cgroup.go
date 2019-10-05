// +build linux
// +build go1.12

package sandbox

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	cgCPUSetPathPrefix = "/sys/fs/cgroup/cpuset/"
	cgCPUPathPrefix    = "/sys/fs/cgroup/cpu/"
	cgPidPathPrefix    = "/sys/fs/cgroup/pids/"
	cgMemoryPathPrefix = "/sys/fs/cgroup/memory/"
)

//noinspection GoUnusedExportedFunction
func InitCGroup(pid, containerID, memory string, cpus string) error {
	_, _ = os.Stderr.WriteString(fmt.Sprintf("DEBUG: InitCGroup(%s, %s, %s) starting...\n", pid, containerID, memory))

	dirs := []string{
		filepath.Join(cgCPUSetPathPrefix, containerID),
		filepath.Join(cgCPUPathPrefix, containerID),
		filepath.Join(cgPidPathPrefix, containerID),
		filepath.Join(cgMemoryPathPrefix, containerID),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("DEBUG: os.MkdirAll(%s, os.ModePerm) failed, err: %s\n", dir, err.Error()))
			return err
		}
	}

	if err := cpusetCGroup(pid, containerID, cpus); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("DEBUG: cpusetCGroup(%s, %s, %s) failed, err: %s\n", pid, containerID, cpus, err.Error()))
		return err
	}

	if err := cpuCGroup(pid, containerID); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("DEBUG: cpuCGroup(%s, %s) failed, err: %s\n", pid, containerID, err.Error()))
		return err
	}

	if err := pidCGroup(pid, containerID); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("DEBUG: pidCGroup(%s, %s) failed, err: %s\n", pid, containerID, err.Error()))
		return err
	}

	if err := memoryCGroup(pid, containerID, memory); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("DEBUG: memoryCGroup(%s, %s) failed, err: %s\n", pid, containerID, err.Error()))
		return err
	}

	_, _ = os.Stderr.WriteString(fmt.Sprintf("DEBUG: InitCGroup(%s, %s, %s) done\n", pid, containerID, memory))
	return nil
}

// https://www.kernel.org/doc/Documentation/cgroup-v1/cpusets.txt
func cpusetCGroup(pid, containerID string, cpus string) error {
	cgCPUsetPath := filepath.Join(cgCPUSetPathPrefix, containerID)
	mapping := map[string]string{
		"cpuset.mems": "0",
		"cpuset.cpus": cpus,
		"tasks":       pid,
	}

	for key, value := range mapping {
		path := filepath.Join(cgCPUsetPath, key)
		if err := ioutil.WriteFile(path, []byte(value), 0644); err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("Writing [%s] to file: %s failed\n", value, path))
			return err
		}
		c, _ := ioutil.ReadFile(path)
		_, _ = os.Stderr.WriteString(fmt.Sprintf("DEBUG: Content of %s is: %s", path, c))
	}
	return nil
}

// https://www.kernel.org/doc/Documentation/scheduler/sched-bwc.txt
func cpuCGroup(pid, containerID string) error {
	cgCPUPath := filepath.Join(cgCPUPathPrefix, containerID)
	mapping := map[string]string{
		"tasks":            pid,
		"cpu.cfs_quota_us": "10000",
	}

	for key, value := range mapping {
		path := filepath.Join(cgCPUPath, key)
		if err := ioutil.WriteFile(path, []byte(value), 0644); err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("Writing [%s] to file: %s failed\n", value, path))
			return err
		}
		c, _ := ioutil.ReadFile(path)
		_, _ = os.Stderr.WriteString(fmt.Sprintf("DEBUG: Content of %s is: %s", path, c))
	}
	return nil
}

// https://www.kernel.org/doc/Documentation/cgroup-v1/pids.txt
func pidCGroup(pid, containerID string) error {
	cgPidPath := filepath.Join(cgPidPathPrefix, containerID)
	mapping := map[string]string{
		"cgroup.procs": pid,
		"pids.max":     "64",
	}

	for key, value := range mapping {
		path := filepath.Join(cgPidPath, key)
		if err := ioutil.WriteFile(path, []byte(value), 0644); err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("Writing [%s] to file: %s failed\n", value, path))
			return err
		}
		c, _ := ioutil.ReadFile(path)
		_, _ = os.Stderr.WriteString(fmt.Sprintf("DEBUG: Content of %s is: %s", path, c))
	}
	return nil
}

// https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt
func memoryCGroup(pid, containerID, memory string) error {
	cgMemoryPath := filepath.Join(cgMemoryPathPrefix, containerID)
	mapping := map[string]string{
		"memory.kmem.limit_in_bytes": "64m",
		"tasks":                      pid,
		"memory.limit_in_bytes":      fmt.Sprintf("%sK", memory),
	}

	for key, value := range mapping {
		path := filepath.Join(cgMemoryPath, key)
		if err := ioutil.WriteFile(path, []byte(value), 0644); err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("Writing [%s] to file: %s failed\n", value, path))
			return err
		}
		c, _ := ioutil.ReadFile(path)
		_, _ = os.Stderr.WriteString(fmt.Sprintf("DEBUG: Content of %s is: %s", path, c))
	}
	return nil
}
