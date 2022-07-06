package lib

import (
	"os"
	"path/filepath"
)

func SearchFiles(root string, ext string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func isPath(callGraph map[FuncDescriptor][]FuncDescriptor, current FuncDescriptor, goal FuncDescriptor, visited map[FuncDescriptor]bool) bool {
	if current == goal {
		return true
	}

	value, ok := callGraph[current]
	if ok {
		for _, child := range value {
			exists := visited[child]
			if exists {
				continue
			}
			visited[child] = true
			if isPath(callGraph, child, goal, visited) {
				return true
			}
		}
	}
	return false
}

func Contains(a []FuncDescriptor, x FuncDescriptor) bool {
	for _, n := range a {
		if x.TypeHash() == n.TypeHash() {
			return true
		}
	}
	return false
}

func Revert(path string) {
	goExt := ".go"
	originalExt := ".original"
	files := SearchFiles(path, ".tmp")
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
