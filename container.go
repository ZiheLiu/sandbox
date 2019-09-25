// +build linux
// +build go1.12

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/ZiheLiu/sandbox/sandbox"
	"github.com/docker/docker/pkg/reexec"
	"github.com/satori/go.uuid"
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
		os.Exit(-1)
	}

	cmd := exec.Command(command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Env = []string{"PS1=[justice] # "}

	time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	})

	startTime := time.Now().UnixNano() / 1e6
	if err := cmd.Run(); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(-1)
	}
	endTime := time.Now().UnixNano() / 1e6

	timeCost, memoryCost := endTime-startTime, cmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss/1024
	_, _ = os.Stderr.WriteString(fmt.Sprintf("INFO: timeCost:%v\n", timeCost))
	_, _ = os.Stderr.WriteString(fmt.Sprintf("INFO: memoryCost:%v\n", memoryCost))
}

// logs will be printed to os.Stderr
func main() {
	basedir := flag.String("basedir", "/tmp", "basedir of tmp binary")
	command := flag.String("command", "./Main", "the command needed to be execute in sandbox")
	timeout := flag.String("timeout", "2000", "timeout in milliseconds")
	memory := flag.String("memory", "256", "memory limitation in MB")
	flag.Parse()

	u := uuid.NewV4()
	if err := sandbox.InitCGroup(strconv.Itoa(os.Getpid()), u.String(), *memory); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(-1)
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
		_, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
		os.Exit(-1)
	}

	os.Exit(0)
}
