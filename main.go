package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
)

func usage() {
	fmt.Println("\nusage autotel --command [path to go project]")
	fmt.Println("\tcommand:")
	fmt.Println("\t\tinject          (injects open telemetry calls into project code)")
	fmt.Println("\t\tdumpcfg         (dumps control flow graph)")
	fmt.Println("\t\tgencfg          (generates json representation of control flow graph)")
	fmt.Println("\t\trootfunctions   (dumps root functions)")
}

func isPath(callGraph map[string]string, current string, goal string) bool {
	if current == goal {
		return true
	}
	value, ok := callGraph[current]
	if ok {
		if isPath(callGraph, value, goal) {
			return true
		}
	}
	return false
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

func findRootFunctions(file string) []string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
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

func buildCallGraph(file string) map[string]string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	currentFun := "nil"
	backwardCallGraph := make(map[string]string)

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			id, ok := x.Fun.(*ast.Ident)
			if ok {
				backwardCallGraph[id.Name] = currentFun
			}
		case *ast.FuncDecl:
			currentFun = x.Name.Name
		}
		return true
	})

	return backwardCallGraph
}

func parsePath(root string) {
	files := searchFiles(root)

	var rootFunctions []string
	backwardCallGraph := make(map[string]string)

	for _, file := range files {
		rootFunctions = append(rootFunctions, findRootFunctions(file)...)
	}
	for _, file := range files {
		callGraphInstance := buildCallGraph(file)
		for key, value := range callGraphInstance {
			backwardCallGraph[key] = value
		}
	}
	for _, file := range files {
		instrument(file, backwardCallGraph, rootFunctions)
	}
}

// var callgraph = {
//     nodes: [
//         { data: { id: 'fun1' } },
//         { data: { id: 'fun2' } },
// 		],
//     edges: [
//         { data: { id: 'e1', source: 'fun1', target: 'fun2' } },
//     ]
// };

func generatecfg(callgraph map[string]string) {
	functions := make(map[string]bool, 0)
	for k, v := range callgraph {
		if functions[k] == false {
			functions[k] = true
		}
		if functions[v] == false {
			functions[v] = true
		}
	}
	for f := range functions {
		fmt.Println(f)
	}
	out, err := os.Create("./ui/callgraph.js")
	defer out.Close()
	if err != nil {
		return
	}
	out.WriteString("var callgraph = {")
	out.WriteString("\n\tnodes: [")
	for f := range functions {
		out.WriteString("\n\t\t { data: { id: '")
		out.WriteString(f)
		out.WriteString("' } },")
	}
	out.WriteString("\n\t],")
	out.WriteString("\n\tedges: [")
	edgeCounter := 0
	for k, v := range callgraph {
		out.WriteString("\n\t\t { data: { id: '")
		out.WriteString("e" + strconv.Itoa(edgeCounter))
		out.WriteString("', ")
		out.WriteString("source: '")
		out.WriteString(v)
		out.WriteString("', ")
		out.WriteString("target: '")
		out.WriteString(k)
		out.WriteString("' ")
		out.WriteString("} },")
		edgeCounter++
	}
	out.WriteString("\n\t]")
	out.WriteString("\n};")
}

// Parsing algorithm works as follows. It goes through all function
// decls and infer function bodies to find call to SumoAutoInstrument
// A parent function of this call will become root of instrumentation
// Each function call from this place will be instrumented automatically

func main() {
	fmt.Println("autotel compiler")
	args := len(os.Args)
	if args < 3 {
		usage()
		return
	}
	if os.Args[1] == "--inject" {
		parsePath(os.Args[2])
		fmt.Println("\tinstrumentation done")
	}
	if os.Args[1] == "--dumpcfg" {
		files := searchFiles(os.Args[2])
		backwardCallGraph := make(map[string]string)
		for _, file := range files {
			callGraphInstance := buildCallGraph(file)
			for key, value := range callGraphInstance {
				backwardCallGraph[key] = value
			}
		}
		for k, v := range backwardCallGraph {
			fmt.Print("\n\t", k)
			fmt.Print(" -> ", v)
		}
	}
	if os.Args[1] == "--gencfg" {
		files := searchFiles(os.Args[2])
		backwardCallGraph := make(map[string]string)
		for _, file := range files {
			callGraphInstance := buildCallGraph(file)
			for key, value := range callGraphInstance {
				backwardCallGraph[key] = value
			}
		}
		generatecfg(backwardCallGraph)
	}
	if os.Args[1] == "--rootfunctions" {
		files := searchFiles(os.Args[2])
		var rootFunctions []string
		for _, file := range files {
			rootFunctions = append(rootFunctions, findRootFunctions(file)...)
		}
		for _, fun := range rootFunctions {
			fmt.Println("\t" + fun)
		}
	}

}
