// +build linux
// +build go1.12

package main

import (
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ZiheLiu/sandbox/sandbox"
	"github.com/docker/docker/pkg/reexec"
)

func init() {
	// register "justiceInit" => justiceInit() every time
	reexec.Register("justiceInit", justiceInit)

	/**
	* 0. `init()` adds key "justiceInit" in `map`;
	* 1. reexec.Init() seeks if key `os.Args[0]` exists in `registeredInitializers`;
	* 2. for the first time this binary is invoked, the key is os.Args[0], AKA "/path/to/clike_container",
	     which `registeredInitializers` will return `false`;
	* 3. `main()` calls binary itself by reexec.Command("justiceInit", args...);
	* 4. for the second time this binary is invoked, the key is os.Args[0], AKA "justiceInit",
	*    which exists in `registeredInitializers`;
	* 5. the value `justiceInit()` is invoked, any hooks(like set hostname) before fork() can be placed here.
	*/
	if reexec.Init() {
		os.Exit(0)
	}
}

func justiceInit() {
	basedir := os.Args[1]
	command := os.Args[2]
	timeout, _ := strconv.ParseInt(os.Args[3], 10, 32)

	if err := sandbox.InitNamespace(basedir); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(0)
	}

	cmd := exec.Command(command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = []string{"PS1=[justice] # "}

	tle := false
	time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
		tle = true
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	})

	startTime := time.Now().UnixNano() / 1e6
	if err := cmd.Run(); err != nil {
		if tle {
			_, _ = os.Stderr.WriteString(fmt.Sprintln("Time Limit Error"))
		} else {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		}
		os.Exit(0)
	}
	endTime := time.Now().UnixNano() / 1e6

	timeCost, memoryCost := endTime-startTime, cmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss/1024
	_, _ = os.Stderr.WriteString(fmt.Sprintf("INFO: timeCost:%v\n", timeCost))
	_, _ = os.Stderr.WriteString(fmt.Sprintf("INFO: memoryCost:%v\n", memoryCost))
}

func cgroupOomControl(containerId string) map[string]string {
	res := make(map[string]string)

	filePath := filepath.Join("/sys/fs/cgroup/memory/", containerId, "memory.oom_control")
	c, _ := ioutil.ReadFile(filePath)
	rows := strings.Split(string(c), "\n")
	for _, row := range rows {
		if row != "" {
			params := strings.Split(row, " ")
			res[params[0]] = params[1]
		}
	}

	return res
}

// logs will be printed to os.Stderr
func main() {
	basedir := flag.String("basedir", "/tmp", "basedir of tmp binary")
	command := flag.String("command", "./Main", "the command needed to be execute in sandbox")
	timeout := flag.String("timeout", "2000", "timeout in milliseconds")
	memory := flag.String("memory", "256", "memory limitation in MB")
	flag.Parse()

	containerId := uuid.NewV4().String()

	if err := sandbox.InitCGroup(strconv.Itoa(os.Getpid()), containerId, *memory); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(0)
	}

	cmd := reexec.Command("justiceInit", *basedir, *command, *timeout)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	if err := cmd.Run(); err != nil {
		oomInfo := cgroupOomControl(containerId)
		if val, ok := oomInfo["oom_kill"]; ok && val != "0" {
			_, _ = os.Stderr.WriteString(fmt.Sprintln("Memory Limit Error"))
		} else {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		}
		os.Exit(0)
	}

	os.Exit(0)
}
