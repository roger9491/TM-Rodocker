package container

import (
	"TM-Rodocker/subsystems"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	fmt.Println(args)
	cmd := exec.Command("/proc/self/exe", args...) // 使用當前的進程重新執行 args ，效果為建立args的子進程，但這個進程是隔離的
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}

func Run(tty bool, command string, res *subsystems.ResourceConfig) {
	parent := NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}
