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
	CPPBaseDir    string
	CPPProjectDir string
)

// copy test source file `*.cpp` to tmp dir
func copyCPPSourceFile(name string, t *testing.T) {
	t.Logf("Copying file %s ...", name)
	if err := os.MkdirAll(CPPBaseDir, os.ModePerm); err != nil {
		t.Errorf("Invoke mkdir(%s) err: %v", CPPBaseDir, err.Error())
	}

	args := []string{
		CPPProjectDir + "/resources/cpp/" + name,
		CPPBaseDir + "/Main.cpp",
	}
	cmd := exec.Command("cp", args...)
	if err := cmd.Run(); err != nil {
		t.Errorf("Invoke `cp %s` err: %v", strings.Join(args, " "), err)
	}
}

// compile CPP source file
func compileCPP(name, baseDir string, t *testing.T) string {
	t.Logf("Compiling file %s ...", name)

	var stderr bytes.Buffer
	args := []string{
		"-compiler=/usr/bin/g++",
		"-basedir=" + baseDir,
		"-filename=Main.cpp",
		"-timeout=3000",
		"-std=gnu++14",
	}
	cmd := exec.Command("/opt/justice-sandbox/bin/clike_compiler", args...)
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Errorf("Invoke `/opt/justice-sandbox/bin/clike_compiler %s` err: %v", strings.Join(args, " "), err)
	}

	return stderr.String()
}

// run binary in our container
func runCPP(baseDir, memory, timeout string, t *testing.T) (string, string) {
	t.Log("Running file /Main ...")

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

	t.Logf("stderr of runCPP: %s", stderr.String())
	return stdout.String(), stderr.String()
}

func TestCPP0000Fixture(t *testing.T) {
	CPPProjectDir, _ = os.Getwd()
	CPPBaseDir = CPPProjectDir + "/tmp"
}

func TestCPP0001AC(t *testing.T) {
	name := "ac.cpp"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCPPSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CPPBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CPPBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileCPP(name, CPPBaseDir, t), ShouldBeEmpty)
		stdout, _ := runCPP(CPPBaseDir, "16000", "1000", t)
		So(stdout, ShouldEqual, "10:10:23\n")
	})
}

func TestCPP0002CompilerBomb1(t *testing.T) {
	name := "compiler_bomb_1.cpp"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCPPSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CPPBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CPPBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileCPP(name, CPPBaseDir, t), ShouldContainSubstring, "signal: killed")
	})
}

func TestCPP0003CompilerBomb2(t *testing.T) {
	name := "compiler_bomb_2.cpp"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCPPSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CPPBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CPPBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileCPP(name, CPPBaseDir, t), ShouldContainSubstring, "compilation terminated due to -fmax-errors=")
	})
}

func TestCPP0004CompilerBomb3(t *testing.T) {
	name := "compiler_bomb_3.cpp"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCPPSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CPPBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CPPBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileCPP(name, CPPBaseDir, t), ShouldContainSubstring, "template instantiation depth exceeds maximum of")
	})
}

func TestCPP0005CompilerBomb4(t *testing.T) {
	name := "compiler_bomb_4.cpp"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCPPSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CPPBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CPPBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileCPP(name, CPPBaseDir, t), ShouldContainSubstring, "signal: killed")
	})
}

func TestCPP0006CoreDump0(t *testing.T) {
	name := "core_dump_0.cpp"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCPPSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CPPBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CPPBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileCPP(name, CPPBaseDir, t), ShouldBeEmpty)
		// terminate called after throwing an instance of 'char const*'
		_, stderr := runCPP(CPPBaseDir, "64000", "1000", t)
		So(stderr, ShouldContainSubstring, "signal: aborted (core dumped)")
	})
}

func TestCPP0007ForkBomb(t *testing.T) {
	name := "fork_bomb.cpp"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCPPSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CPPBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CPPBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileCPP(name, CPPBaseDir, t), ShouldBeEmpty)
		_, stderr := runCPP(CPPBaseDir, "64000", "1000", t)
		So(stderr, ShouldContainSubstring, "Time Limit Error")
	})
}

func TestCPP0008IncludeLeaks(t *testing.T) {
	name := "include_leaks.cpp"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCPPSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CPPBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CPPBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileCPP(name, CPPBaseDir, t), ShouldContainSubstring, "/etc/shadow")
	})
}

func TestCPP0009InfiniteLoop(t *testing.T) {
	name := "infinite_loop.cpp"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCPPSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CPPBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CPPBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileCPP(name, CPPBaseDir, t), ShouldBeEmpty)
		_, stderr := runCPP(CPPBaseDir, "64000", "1000", t)
		So(stderr, ShouldContainSubstring, "Time Limit Error")
	})
}

func TestCPP0010MemoryAllocation(t *testing.T) {
	name := "memory_allocation.cpp"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCPPSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CPPBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CPPBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileCPP(name, CPPBaseDir, t), ShouldBeEmpty)
		_, stderr := runCPP(CPPBaseDir, "800", "1000", t)
		So(stderr, ShouldContainSubstring, "Memory Limit Error")
	})
}

func TestCPP0011PlainText(t *testing.T) {
	name := "plain_text.cpp"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCPPSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CPPBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CPPBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileCPP(name, CPPBaseDir, t), ShouldContainSubstring, "error")
	})
}

func TestCPP0012RunCommandLine0(t *testing.T) {
	name := "run_command_line_0.cpp"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCPPSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CPPBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CPPBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileCPP(name, CPPBaseDir, t), ShouldBeEmpty)
		stdin, _ := runCPP(CPPBaseDir, "16000", "1000", t)
		So(stdin, ShouldEqual, "32512 32512")
	})
}

func TestCPP0013RunWriteFile(t *testing.T) {
	name := "write_file.cpp"
	Convey(fmt.Sprintf("Testing [%s]...", name), t, func() {
		copyCPPSourceFile(name, t)
		defer func() {
			if err := os.RemoveAll(CPPBaseDir); err != nil {
				t.Errorf("Invoke `os.RemoveAll(%s)` err: %v", CPPBaseDir, err)
				t.FailNow()
			}
		}()

		So(compileCPP(name, CPPBaseDir, t), ShouldBeEmpty)
		stdin, _ := runCPP(CPPBaseDir, "16000", "1000", t)
		So(stdin, ShouldEqual, "write to file /test.txt failed\n")
	})
}
