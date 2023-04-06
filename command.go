package main

import (
	"TM-Rodocker/cgroups/subsystems"
	"TM-Rodocker/container"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run", // 命令名稱
	Usage: `Create a  container with namespace and cgroups limit
			my docker run -ti [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{ // 有輸入 ti 就是 true
			Name:  "ti",         // run -ti
			Usage: "enable tty", // 顯示功能
		},
		cli.StringFlag{ // -m/ --m : 是一個帶字串參數的選項
			Name:  "m",                   // run -ti
			Usage: "enable memory limit", // 顯示功能
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
	},

	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			err := errors.New("An error occurred")
			log.Error(err)

			// 返回 error 结构
			return err
		}

		// 命令參數串列
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		fmt.Println("cmdArray is ", cmdArray)
		

		resConf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"), // 從 指標m 獲取值
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}

		tty := context.Bool("ti") // 是否輸入 -ti

		volume := context.String("v")
		fmt.Println("volume is ", volume)

		Run(tty, cmdArray, resConf, volume) // 執行命令
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
