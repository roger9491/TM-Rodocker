package main

import (
	"TM-Rodocker/cgroups"
	"TM-Rodocker/cgroups/subsystems"
	"TM-Rodocker/container"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Run(tty bool, cmdArray []string, res *subsystems.ResourceConfig, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// 創建 cgroup 實例
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")

	// 結束時刪除 cgroup
	defer cgroupManager.Destroy()
	// 設定 cgroup 資源限制
	cgroupManager.Set(res)
	// 將進程加入到 cgroup 中
	cgroupManager.Apply(parent.Process.Pid)

	// 發送用戶命令
	sendInitCommand(cmdArray, writePipe)

	parent.Wait()

	// 刪除工作層
	mntURL := "/root/mnt/"
	rootURL := "/root/"

	container.DeleteWorkSpace(rootURL, mntURL, volume)
	os.Exit(-1)
}

func sendInitCommand(cmdArray []string, writePipe *os.File) {
	command := strings.Join(cmdArray, " ")
	log.Infof("command all is %s", command)

	writePipe.WriteString(command)
	writePipe.Close()
}
