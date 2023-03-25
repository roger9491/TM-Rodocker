package container

import (
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func RunContainerInitProcess(command string, args []string) error {
	log.Infof("command %s", command)

	// proc 是一個偽文件系統，主要用於訪問內核和進程相關訊息 ex. ps
	// defaultMountFlags: 設置標誌
	// MS_NOEXEC: 在此挂载点上不允许执行二进制文件。
	// MS_NOSUID: 在此挂载点上不允许 set-user-ID 和 set-group-ID。
	// MS_NODEV: 在此挂载点上不允许访问设备文件。
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	argv := []string{command}

	// syscall.Exec 啟動新的進程，替換當前進程的
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}

	return nil
}
