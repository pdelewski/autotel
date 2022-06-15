package autotellib

import (
	"os"
)

const (
	contextPassFileSuffix         = "_pass_ctx"
	instrumentationPassFileSuffix = "_pass_tracing.go"
)

func ExecutePasses(files []string, rootFunctions []string, backwardCallGraph map[string][]string) {
	funcDecls := FindCompleteFuncDecls(files)

	for _, file := range files {
		PropagateContext(file, backwardCallGraph, rootFunctions, funcDecls, contextPassFileSuffix)
	}
	for _, file := range files {
		Instrument(file+contextPassFileSuffix, backwardCallGraph, rootFunctions, instrumentationPassFileSuffix)
	}
	for _, file := range files {
		os.Rename(file, file+".original")
	}
}
