package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func usage() {
	fmt.Println("\nusage autotel [path to go project]")
}

func searchFiles(root string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func instrument(file string) {
	fmt.Println("instrumentation", file)
}

func parse(root string) {
	fmt.Println("parsing", root)
	files := searchFiles(root)
	for _, file := range files {
		fmt.Println("pass 1", file)
	}
	for _, file := range files {
		instrument(file)
	}
}

func main() {
	fmt.Println("autotel compiler")
	args := len(os.Args)
	if args != 2 {
		usage()
		return
	}
	parse(os.Args[1])

}
