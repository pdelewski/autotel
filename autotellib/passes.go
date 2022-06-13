package autotellib

import (
	"os"
)

const (
	contextPassFileSuffix         = "_pass_ctx"
	instrumentationPassFileSuffix = "_pass_tracing.go"
)

func ExecutePasses(files []string, rootFunctions []string, backwardCallGraph map[string][]string) {
	funcDecls := make(map[string]bool)
	for _, file := range files {
		funcDeclsFile := FindFuncDecls(file)
		for k, v := range funcDeclsFile {
			funcDecls[k] = v
		}
	}

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

func Revert(path string) {
	goExt := ".go"
	originalExt := ".original"
	files := SearchFiles(path, goExt+contextPassFileSuffix)
	for _, file := range files {
		os.Remove(file)
	}
	files = SearchFiles(path, goExt)
	for _, file := range files {
		os.Remove(file)
	}
	files = SearchFiles(path, originalExt)
	for _, file := range files {
		newFile := file[:len(file)-(len(goExt)+len(originalExt))]
		os.Rename(file, newFile+".go")
	}
}
