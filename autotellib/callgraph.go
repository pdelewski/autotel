package autotellib

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
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

func FindRootFunctions(file string) []string {
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

func BuildCallGraph(file string) map[string][]string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	currentFun := "nil"
	backwardCallGraph := make(map[string][]string)

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			id, ok := x.Fun.(*ast.Ident)
			if ok {
				backwardCallGraph[id.Name] = append(backwardCallGraph[id.Name], currentFun)
			}
			sel, ok := x.Fun.(*ast.SelectorExpr)
			if ok {
				backwardCallGraph[sel.Sel.Name] = append(backwardCallGraph[sel.Sel.Name], currentFun)
			}
		case *ast.FuncDecl:
			currentFun = x.Name.Name
		}
		return true
	})

	return backwardCallGraph
}

func FindFuncDecls(file string) map[string]bool {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	funcDecls := make(map[string]bool)

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			funcDecls[x.Name.Name] = true
		}
		return true
	})

	return funcDecls
}

func InferRootFunctionsFromGraph(callgraph map[string][]string) []string {
	var allFunctions map[string]bool
	var rootFunctions []string
	allFunctions = make(map[string]bool)
	for k, v := range callgraph {
		allFunctions[k] = true
		for _, childFun := range v {
			allFunctions[childFun] = true
		}
	}
	for k, _ := range allFunctions {
		_, exists := callgraph[k]
		if !exists {
			rootFunctions = append(rootFunctions, k)
		}
	}
	return rootFunctions
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

func Generatecfg(callgraph map[string][]string, path string) {
	functions := make(map[string]bool, 0)
	for k, childFuns := range callgraph {
		if functions[k] == false {
			functions[k] = true
		}
		for _, v := range childFuns {
			if functions[v] == false {
				functions[v] = true
			}
		}
	}
	for f := range functions {
		fmt.Println(f)
	}
	out, err := os.Create(path)
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
	for k, children := range callgraph {
		for _, childFun := range children {
			out.WriteString("\n\t\t { data: { id: '")
			out.WriteString("e" + strconv.Itoa(edgeCounter))
			out.WriteString("', ")
			out.WriteString("source: '")

			out.WriteString(childFun)
			out.WriteString(" ")

			out.WriteString("', ")
			out.WriteString("target: '")
			out.WriteString(k)
			out.WriteString("' ")
			out.WriteString("} },")
			edgeCounter++
		}
	}
	out.WriteString("\n\t]")
	out.WriteString("\n};")
}
