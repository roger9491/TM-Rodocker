package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = "usage"

	// 提供的命令列表
	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	// 在命令执行前，设置日志格式
	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})

		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
