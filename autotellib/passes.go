package autotellib

import (
	"os"
)

const (
	contextPassFileSuffix         = "_pass_ctx"
	instrumentationPassFileSuffix = "_pass_tracing.go"
)

func ExecutePasses(projectPath string, packagePattern string, files []string, rootFunctions []string, funcDecls map[string]bool, backwardCallGraph map[string][]string) {

	GlobalPropagateContext(projectPath, packagePattern, backwardCallGraph, rootFunctions, funcDecls, contextPassFileSuffix)
	for _, file := range files {
		Instrument(file+contextPassFileSuffix, backwardCallGraph, rootFunctions, instrumentationPassFileSuffix)
	}
	for _, file := range files {
		os.Rename(file, file+".original")
	}
}
