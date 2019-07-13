package main

import (
	"os"

	logs "github.com/appscode/go/log/golog"
	"pharmer.dev/cloud-controller-manager/cmds"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewRootCmd(Version).Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
