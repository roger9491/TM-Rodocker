package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

/*
	建立與宿主機隔離的文件系統
	使得容器內的文件操作不會影響到宿主機
*/
// 創建容器的文件系統
func NewWorkSpace(rootURL string, mntURL string, volume string) {
	// 建立只讀層
	CreateReadOnlyLayer(rootURL)
	// 建立可寫層
	CreateWriteLayer(rootURL)
	// 建立挂載點
	CreateMountPoint(rootURL, mntURL)

	// 根據 volume 判斷是否執行掛載數據捲操作
	if volume != "" {
		volumeURLs, ok := volumeUrlExtract(volume)
		fmt.Println("ok: ", ok)
		if ok {
			
			MountVolume(rootURL, mntURL, volumeURLs)
			log.Info("%q", volumeURLs)
		}else{
			log.Infof("Volume parameter input error")
		}
	}
}

// 建立只讀層
func CreateReadOnlyLayer(rootURL string) {
	busyboxURL := rootURL + "busybox/"
	busyboxTarURL := rootURL + "busybox.tar"
	exist, err := PathExists(busyboxURL)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists. %v", busyboxURL, err)
	}

	// 不存在資料夾
	if exist == false {
		// 建立資料夾
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			log.Errorf("Mkdir dir %s error %v", busyboxURL, err)
		}

		// 解壓縮 busybox.tar 到 busyboxURL
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error %v", busyboxURL, err)
		}
	}
}

// 創建 writeLayer 文件夾作為容器唯一的可寫層
func CreateWriteLayer(rootURL string) {
	writeURL := rootURL + WriteLayerName
	if err := os.Mkdir(writeURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error %v", writeURL, err)
	}
}

// 創建堆疊文件系統
func CreateMountPoint(rootURL string, mntURL string) {
	// 創建 mnt 文件夾作為掛載點，此文件為用戶操作的目錄，也是merge層
	if err := os.Mkdir(mntURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error %v", mntURL, err)
	}

	// 創建 overlayfs 所需的 work 層
	workDir := rootURL + WorkerName
	if err := os.Mkdir(workDir, 0777); err != nil {
		log.Errorf("Mkdir dir %s error %v", workDir, err)
	}

	/*
		mount -t overlay overlay -o lowerdir=<t1:t2:...>, upperdir=<u>, workdir=<w> <merge>
		mount命令用于挂载文件系统，
		-t选项后面跟的是文件系统类型，此例中为overlay。
		-o选项后面跟的是挂载选项，
		lowerdir: 指定只读层的路径，可以有多个只读层，用冒号分隔。
		upperdir: 指定可读写层的路径。
		workdir: 指定work目录的路径，用於暫存 overlayfs 臨時文件。
		<merge>: 指定 merge 目录的路径，表示最終 lower 和 upper 的合併目录。
	*/
	lowDir := "lowerdir=" + rootURL + BusyboxName
	upperDir := "upperdir=" + rootURL + WriteLayerName
	workDir = "workdir=" + workDir
	dirs := lowDir + "," + upperDir + "," + workDir
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Mount point error %v", err)
	}

}

// 判斷文件路徑是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/*
刪除容器的文件系統
*/
func DeleteWorkSpace(rootURL string, mntURL string, volume string) {
	// 根據 volume 判斷是否執行卸載數據捲操作
	if volume != "" {
		volumeURLs, ok := volumeUrlExtract(volume)
		if ok {
			UmountVolume(mntURL, volumeURLs)
		}
	}
	DeleteMountPoint(rootURL, mntURL)
	DeleteWriteLayer(rootURL)

}

func DeleteMountPoint(rootURL string, mntURL string) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Unmount dir %s error %v", mntURL, err)
	}

	if err := os.RemoveAll(mntURL); err != nil {
		log.Errorf("Remove dir %s error %v", mntURL, err)
	}
}

func DeleteWriteLayer(rootURL string) {
	writeURL := rootURL + "writeLayer"
	if err := os.RemoveAll(writeURL); err != nil {
		log.Errorf("Remove dir %s error %v", writeURL, err)
	}
}

// 卸載數據捲
func UmountVolume(mntURL string, volumeURLs []string) {
	containerUrl := mntURL + volumeURLs[1]
	cmd := exec.Command("umount", containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil { 
		log.Errorf("Unmount volume faild. %v", err)
	}
}

// 解析volume參數，並檢驗
func volumeUrlExtract(volume string) (volumeURLs []string, ok bool) {
	ok = false
	volumeURLs = strings.Split(volume, ":")
	if len(volumeURLs) == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
		ok = true
		return
	}
	return
}

// 掛載數據捲
func MountVolume(rootURL string, mntURL string, volumeURLs []string) {
	// 創建宿主機文件目錄
	parentUrl := volumeURLs[0]
	if err := os.Mkdir(parentUrl, 0777); err != nil {
		log.Infof("Mkdir dir %s error %v", parentUrl, err)
	}

	// 容器文件系統裡創建掛載目錄
	containerUrl := volumeURLs[1]
	containerVolumeUrl := mntURL + containerUrl
	if err := os.Mkdir(containerVolumeUrl, 0777); err != nil {
		log.Infof("Mkdir container dir %s error %v", containerVolumeUrl, err)
	}

	// 將宿主機文件目錄掛載到容器系統裡
	cmd := exec.Command("mount", "--bind", parentUrl,  containerVolumeUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Mount volum failed %v", err)
	}
}