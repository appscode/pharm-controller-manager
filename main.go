package main

import (
	"os"

	logs "github.com/appscode/log/golog"
	"github.com/appscode/pharm-controller-manager/cmds"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
