package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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

func Parse(file string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		fmt.Println(fn.Name.Name)
	}
	var currentFun string
	ast.Inspect(node, func(n ast.Node) bool {
		// Find Functions
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			currentFun = fn.Name.Name
			return true
		}
		funcCall, ok := n.(*ast.CallExpr)
		if ok {
			fmt.Println("FuncCall: ", funcCall.Fun)
			fmt.Println("child of: ", currentFun)
		}
		return true
	})
}

func findRootFunctions(file string) []string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		fmt.Println(fn.Name.Name)
	}
	var currentFun string
	var rootFunctions []string

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			_, ok := x.Fun.(*ast.Ident)
			if ok {
			}
			selector, ok := x.Fun.(*ast.SelectorExpr)
			if ok {
				if selector.Sel.Name == "SumoAutoInstrument" {
					rootFunctions = append(rootFunctions, currentFun)
				}
			}
		case *ast.FuncDecl:
			currentFun = x.Name.Name
		}
		return true
	})

	return rootFunctions
}

func parsePath(root string) {
	fmt.Println("parsing", root)
	files := searchFiles(root)
	for _, file := range files {
		fmt.Println("pass 1", file)
	}
	var rootFunctions []string
	for _, file := range files {
		instrument(file)
		rootFunctions = append(rootFunctions, findRootFunctions(file)...)
	}
	fmt.Println("Root Functions:")
	for _, file := range rootFunctions {
		fmt.Println(file)
	}
}

// Parsing algorithm works as follows. It goes through all function
// decls and infer function bodies to find call to SumoAutoInstrument
// A parent function of this call will become root of instrumentation
// Each function call from this place will be instrumented automatically

func main() {
	fmt.Println("autotel compiler")
	args := len(os.Args)
	if args != 2 {
		usage()
		return
	}
	parsePath(os.Args[1])

}
