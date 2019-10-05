package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	CBaseDir    string
	CProjectDir string
)

// copy test source file `*.c` to tmp dir
func copyCSourceFile(name string, t *testing.T) {
	t.Logf("Copying file %s ...", name)
	if err := os.MkdirAll(CBaseDir, os.ModePerm); err != nil {
		t.Errorf("Invoke mkdir(%s) err: %v", CBaseDir, err.Error())
	}

	args := []string{
		CProjectDir + "/resources/c/" + name,
		CBaseDir + "/Main.c",
	}
	cmd := exec.Command("cp", args...)
	if err := cmd.Run(); err != nil {
		t.Errorf("Invoke `cp %s` err: %v", strings.Join(args, " "), err)
	}
}

// compile C source file
func compileC(name, baseDir string, t *testing.T) string {
	t.Logf("Compiling file %s ...", name)

	var stderr bytes.Buffer
	args := []string{
		"-compiler=/usr/bin/gcc",
		"-basedir=" + baseDir,
		"-filename=Main.c",
		"-timeout=3000",
		"-std=gnu11",
	}
	cmd := exec.Command("/opt/justice-sandbox/bin/clike_compiler", args...)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Errorf("Invoke `/opt/justice-sandbox/bin/clike_compiler %s` err: %v", strings.Join(args, " "), err)
	}

	return stderr.String()
}

// run binary in our container
func runC(baseDir, memory, timeout string, t *testing.T) (string, string) {
	t.Log("Running binary /Main ...")

	var stdout, stderr bytes.Buffer
	args := []string{
		"-basedir=" + baseDir,
		"-memory=" + memory,
		"-timeout=" + timeout,
		"-command=./Main",
		"-username=oj-user",
	}
	cmd := exec.Command("/opt/justice-sandbox/bin/clike_container", args...)
	cmd.Stdin = strings.NewReader("10:10:23AM")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Errorf("Invoke `/opt/justice-sandbox/bin/clike_container %s` err: %v", strings.Join(args, " "), err)
	}

	t.Logf("stderr of runC: %s", stderr.String())
	return stdout.String(), stderr.String()
}

func TestC0000Fixture(t *testing.T) {
	CProjectDir, _ = os.Getwd()
	CBaseDir = CProjectDir + "/tmp"
}

func TestC0001AC(t *testing.T) {
	name := "ac.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		stdout, _ := runC(CBaseDir, "64000", "1000", t)
		So(stdout, ShouldContainSubstring, "10:10:23")
	})
}

func TestC0002CompilerBomb0(t *testing.T) {
	name := "compiler_bomb_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldContainSubstring, "signal: killed")
	})
}

func TestC0003CompilerBomb1(t *testing.T) {
	name := "compiler_bomb_1.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldContainSubstring, "signal: killed")
	})
}

func TestC0004CompilerBomb2(t *testing.T) {
	name := "compiler_bomb_2.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldContainSubstring, "signal: killed")
	})
}

func TestC0005CompilerBomb3(t *testing.T) {
	name := "compiler_bomb_3.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldContainSubstring, "signal: killed")
	})
}

func TestC0006CoreDump0(t *testing.T) {
	name := "core_dump_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		_, stderr := runC(CBaseDir, "64000", "1000", t)
		So(stderr, ShouldContainSubstring, "signal: segmentation fault (core dumped)")
	})
}

func TestC0007CoreDump1(t *testing.T) {
	name := "core_dump_1.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		// warning: division by zero [-Wdiv-by-zero]
		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		_, stderr := runC(CBaseDir, "64000", "1000", t)
		So(stderr, ShouldContainSubstring, "signal: floating point exception (core dumped)")
	})
}

func TestC0008CoreDump2(t *testing.T) {
	name := "core_dump_2.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		// *** stack smashing detected ***: terminated
		_, stderr := runC(CBaseDir, "64000", "1000", t)
		So(stderr, ShouldContainSubstring, "signal: aborted (core dumped)")
	})
}

func TestC0009ForkBomb0(t *testing.T) {
	name := "fork_bomb_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		// got `signal: killed`
		_, stderr := runC(CBaseDir, "64000", "1000", t)
		So(stderr, ShouldContainSubstring, "Time Limit Error")
	})
}

func TestC0010ForkBomb1(t *testing.T) {
	name := "fork_bomb_1.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		// got `signal: killed`
		_, stderr := runC(CBaseDir, "64000", "1000", t)
		So(stderr, ShouldContainSubstring, "signal: segmentation fault (core dumped)")
	})
}

func TestC0011GetHostByName(t *testing.T) {
	name := "get_host_by_name.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		// Main.c:(.text+0x28): warning: Using 'gethostbyname' in statically linked applications
		// requires at runtime the shared libraries from the glibc version used for linking
		// got `exit status 1`
		stdin, _ := runC(CBaseDir, "64000", "1000", t)
		So(stdin, ShouldContainSubstring, "gethostbyname error")
	})
}

func TestC0012IncludeLeaks(t *testing.T) {
	name := "include_leaks.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldContainSubstring, "/etc/shadow")
	})
}

func TestC0013InfiniteLoop(t *testing.T) {
	name := "infinite_loop.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		// got `signal: killed`
		_, stderr := runC(CBaseDir, "64000", "1000", t)
		So(stderr, ShouldContainSubstring, "Time Limit Error")
	})
}

func TestC0014MemoryAllocation(t *testing.T) {
	name := "memory_allocation.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		_, stderr := runC(CBaseDir, "500", "1000", t)
		So(stderr, ShouldContainSubstring, "Memory Limit Error")
	})
}

func TestC0015PlainText(t *testing.T) {
	name := "plain_text.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldContainSubstring, "error")
	})
}

func TestC0016RunCommandLine0(t *testing.T) {
	name := "run_command_line_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		stdin, _ := runC(CBaseDir, "64000", "1000", t)
		So(stdin, ShouldEqual, "32512")
	})
}

func TestC0017RunCommandLine1(t *testing.T) {
	name := "run_command_line_1.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		stdin, _ := runC(CBaseDir, "64000", "1000", t)
		So(stdin, ShouldContainSubstring, "32512")
	})
}

func TestC0018Syscall0(t *testing.T) {
	name := "syscall_0.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		stdin, _ := runC(CBaseDir, "16000", "1000", t)
		So(stdin, ShouldContainSubstring, "-1")
	})
}

func TestC0019TCPClient(t *testing.T) {
	name := "tcp_client.c"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileC(name, CBaseDir, t), ShouldBeEmpty)
		stdin, _ := runC(CBaseDir, "16000", "1000", t)
		So(stdin, ShouldContainSubstring, "connect failed")
	})
}
