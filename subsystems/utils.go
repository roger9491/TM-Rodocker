package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

/*
	/proc/self/mountinfo
	可以找出與當前進程相關的 mount 訊息
*/

// 通過 /proc/self/mountinfo 獲取 subsystem 的 hierarchy cgroyp 所在的目錄
func FindCgroupMountpoint(subsystem string) string {
	// 打開/proc/self/mountinfo文件
	mountinfo, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer mountinfo.Close()

	// 通過bufio.NewScanner 來讀取文件
	scanner := bufio.NewScanner(mountinfo)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		// 9th field is mountpoint
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			// 最後一段文字 == subsystem
			if opt == subsystem {
				return fields[4] // 路徑
			}
		}
	}

	return ""
}

// 獲取 cgroup 所在的路徑
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)

	// os.IsNotExist(err): 	文件、目錄 存在 true
	// 						文件、目錄 不存在 false
	// 獲取目錄訊息，如果目錄不存在，則創建
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err != nil {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}
