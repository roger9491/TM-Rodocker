package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// 初始化容器進程，設置掛載點
func RunContainerInitProcess(command string, args []string) error {

	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}

	log.Infof("command %s", command)
	
	setUpMount()

	// 尋找命令的絕對路徑
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}
	log.Infof("Find path %s", path)

	// syscall.Exec 不會再 PATH裡面搜尋命令，所以要給絕對路徑
	// syscall.Exec 啟動新的進程，替換當前進程的
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}

	return nil
}

// 讀取管道內容，獲取命令
func readUserCommand() []string {
	// 獲取文件描述符3的內容
	pipe := os.NewFile(uintptr(3), "pipe")

	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)

	return strings.Split(msgStr, " ")
}

// 設置掛載點
func setUpMount() {
	// 獲取當前目錄路徑
	pwd, err := os.Getwd()	
	if err != nil {
		log.Errorf("Get current location error %v", err)
		return
	}

	log.Infof("Current location is %s", pwd)
	pivotRoot(pwd)

	// proc 是一個偽文件系統，主要用於訪問內核和進程相關訊息 ex. ps
	// defaultMountFlags: 設置標誌
	// MS_NOEXEC: 在此挂载点上不允许执行二进制文件。
	// MS_NOSUID: 在此挂载点上不允许 set-user-ID 和 set-group-ID。
	// MS_NODEV: 在此挂载点上不允许访问设备文件。
	// mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}



// 切換當前目錄為新的根目錄
func pivotRoot(root string) error {

	// 效果為 獨立出來一個掛載點
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}

	// 創建 rootfs/.pivot_root 存儲 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}
	/*  
		https://man7.org/linux/man-pages/man2/pivot_root.2.html
	關於 pivot_root 的 使用條件說明，
	簡單說 new_root 和 put_old 不能位於當前根目錄的同一掛載點。
	*/
	// pivot_root 會將當前的 rootfs 移動到 rootfs/.pivot_root，並且將 root 作為新的根目錄
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	// 切換工作目錄到根目錄
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	// unmount rootfs/.pivot_root
	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	// 刪除臨時文件夾
	return os.Remove(pivotDir)
}










