package container

import (
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
)

/*
*exec.Cmd: 命令實例
		設置隔離，文件描述符

*os.File: 寫管道
*/
// 配置進程NameSpcae隔離，並配置文件系統
func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
	// 建立匿名管道
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}

	/*

			NEWUTS: UTS 命名空间，允许每个容器拥有独立的主机名和域名。
			NEWPID: PID 命名空间，为每个容器提供独立的进程 ID 空间，
					使容器中的进程 ID 与主机上其他进程 ID 隔离。
			NEWNS: 文件系统命名空间，允许每个容器拥有独立的文件系统视图，隔离文件系统挂载点。
			NEWNET: 网络命名空间，为每个容器提供独立的网络接口、路由和防火墙规则。
		    NEWIPC: IPC 命名空间，使每个容器拥有独立的 System V IPC、POSIX 消息队列和信号量。
	*/
	// 建立父進程
	cmd := exec.Command("/proc/self/exe", "init") // 使用當前的進程重新執行 args ，效果為建立args的子進程，但這個進程是隔離的
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
		// Cloneflags: 創建一個新的進程，並調用這些參數以達到隔離
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// 傳遞 readPipe 文件描述符給子進程
	cmd.ExtraFiles = []*os.File{readPipe}
	mntURL := "/root/mnt/"
	rootURL := "/root/"
	NewWorkSpace(rootURL, mntURL, volume)
	cmd.Dir = mntURL
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {

	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	return read, write, nil
}
