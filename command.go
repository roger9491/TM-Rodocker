package main

import (
	"TM-Rodocker/container"
	"TM-Rodocker/subsystems"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run", // 命令名稱
	Usage: `Create a  container with namespace and cgroups limit
			my docker run -ti [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",         // run -ti
			Usage: "enable tty", // 顯示功能
		},
		cli.BoolFlag{
			Name:  "m",         // run -ti
			Usage: "enable memory limit", // 顯示功能
		},
	},

	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}

		
		// 命令參數串列
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),	// 從 指標m 獲取值
			CpuSet: 	context.String("cpuset"),
			CpuShare: 	context.String("cpushare"),
		}

		tty := context.Bool("ti") // 是否輸入 -ti
		container.Run(tty, cmdArray, resConf)   // 執行命令
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",

	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		cmd := context.Args().Get(0)
		log.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}
