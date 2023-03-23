package main

import (
	"log"
	"os"
	"syscall"

	"github.com/containerd/containerd/defaults"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/cloud/functions/v1"
)

var runCommand = cli.Command {
	Name: "run", 
	Usage: `Create a  container with namespace and cgroups limit
			my docker run -ti [command]`, 
	Flags: []cli.Flag {
		cli.BoolFlag{
			Name: "ti", 
			Usage: "enable tty",
		},
	}

	Action: func (context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("ti")
		Run(tty, cmd)
		return nil
	}, 
}


var initCommand = cli.Command {
	Name: "init",
	Usage: "Init container process run user's process in container. Do not call it outside",

	Action: func (context *cli.Context) error {
		log.Infof("init come on")
		cmd := context.Args().Get(0)
		log.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}

func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)	// 使用 ?
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS
		| syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}

func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err!= nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}

func RunContainerInitProcess(command string, args []string) error { 
	logrus.Infof("command %s", command)


	// proc 是一個偽文件系統，主要用於訪問內核和進程相關訊息 ex. ps
	// defaultMountFlags: 設置標誌
	// MS_NOEXEC: 在此挂载点上不允许执行二进制文件。
	// MS_NOSUID: 在此挂载点上不允许 set-user-ID 和 set-group-ID。
	// MS_NODEV: 在此挂载点上不允许访问设备文件。
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV 
	syscall.Mount("proc", "/proc", "proc", unitptr(defaultMountFlags), "")

	args := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil{
		log.Errorf(err.Error())
	}

	return nil
}