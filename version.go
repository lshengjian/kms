package main

import (
	"fmt"
	"runtime"
	"github.com/lshengjian/kms/util"
)

var cmdVersion = &Command{
	ExecuteFunc:       runVersion,
	UsageLine: "version",
	Short:     "print KMS version",
	Long:      `print the KMS version`,
}

func runVersion(cmd *Command, args []string) bool {
	if len(args) != 0 {
		cmd.Usage()
	}

	fmt.Printf("version %s %s %s\n", util.VERSION, runtime.GOOS, runtime.GOARCH)
	return true
}
