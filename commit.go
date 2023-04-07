package main

import (
	"TM-Rodocker/global"
	"os/exec"
	log "github.com/sirupsen/logrus"
)

// commitContainer 將容器的工作層打包成 tar 包
func commitContainer(imageName string) {
	imageTar := "/root/" + imageName + ".tar"

	log.Infof("imageTar: %s", imageTar)

	/* "-czf"：指定 tar 命令的选项。-c 表示创建新归档文件；-z 表示使用 gzip 压缩；-f 表示指定输出文件名。
	imageTar：指定输出文件名（由 -f 选项使用）。
	"-C"：指定 tar 命令在执行之前更改到的目录。
	global.MntURL：指定要更改到的目录。
	"."：指定要添加到归档文件中的文件和目录。这里使用 . 表示当前目录（即 global.MntURL）中的所有文件和目录。 */
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", global.MntURL, ".").CombinedOutput(); err != nil {
		log.Errorf("Tar folder %s error %v", imageTar, err)
	}

}
