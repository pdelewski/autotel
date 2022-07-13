package lib

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"os"
	"strconv"

	"golang.org/x/tools/go/packages"
)

type FuncDescriptor struct {
	Id       string
	DeclType string
}

func (fd FuncDescriptor) TypeHash() string {
	return fd.Id + fd.DeclType
}

const mode packages.LoadMode = packages.NeedName |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedFiles

func FindRootFunctions(projectPath string, packagePattern string) []FuncDescriptor {
	fset := token.NewFileSet()

	var currentFun FuncDescriptor
	var rootFunctions []FuncDescriptor

	fmt.Println("FindRootFunctions")
	cfg := &packages.Config{Fset: fset, Mode: mode, Dir: projectPath}
	pkgs, err := packages.Load(cfg, packagePattern)
	if err != nil {
		log.Fatal(err)
	}
	for _, pkg := range pkgs {
		fmt.Println("\t", pkg)

		for _, node := range pkg.Syntax {
			fmt.Println("\t\t", fset.File(node.Pos()).Name())
			ast.Inspect(node, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.CallExpr:
					selector, ok := x.Fun.(*ast.SelectorExpr)
					if ok {
						if selector.Sel.Name == "AutotelEntryPoint__" {
							rootFunctions = append(rootFunctions, currentFun)
						}
					}
				case *ast.FuncDecl:
					currentFun = FuncDescriptor{pkg.TypesInfo.Defs[x.Name].Id(), pkg.TypesInfo.Defs[x.Name].Type().String()}
					fmt.Println("\t\t\tFuncDecl:", pkg.TypesInfo.Defs[x.Name].Id(), pkg.TypesInfo.Defs[x.Name].Type().String())
				}
				return true
			})
		}
	}
	return rootFunctions
}

func BuildCallGraph(projectPath string, packagePattern string, funcDecls map[FuncDescriptor]bool) map[FuncDescriptor][]FuncDescriptor {
	fset := token.NewFileSet()
	cfg := &packages.Config{Fset: fset, Mode: mode, Dir: projectPath}
	pkgs, err := packages.Load(cfg, packagePattern)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("BuildCallGraph")
	currentFun := FuncDescriptor{"nil", ""}
	backwardCallGraph := make(map[FuncDescriptor][]FuncDescriptor)
	for _, pkg := range pkgs {
		fmt.Println("\t", pkg)
		for _, node := range pkg.Syntax {
			fmt.Println("\t\t", fset.File(node.Pos()).Name())
			ast.Inspect(node, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.CallExpr:
					id, ok := x.Fun.(*ast.Ident)
					if ok {
						fmt.Println("\t\t\tFuncCall:", pkg.TypesInfo.Uses[id].Id(), pkg.TypesInfo.Uses[id].Type().String())
						fun := FuncDescriptor{pkg.TypesInfo.Uses[id].Id(), pkg.TypesInfo.Uses[id].Type().String()}
						if !Contains(backwardCallGraph[fun], currentFun) {
							if funcDecls[fun] == true {
								backwardCallGraph[fun] = append(backwardCallGraph[fun], currentFun)
							}
						}
					}
					sel, ok := x.Fun.(*ast.SelectorExpr)
					if ok {
						if pkg.TypesInfo.Uses[sel.Sel] != nil {
							fmt.Println("\t\t\tFuncCall via selector:", pkg.TypesInfo.Uses[sel.Sel].Id(), pkg.TypesInfo.Uses[sel.Sel].Type().String())
							fun := FuncDescriptor{pkg.TypesInfo.Uses[sel.Sel].Id(), pkg.TypesInfo.Uses[sel.Sel].Type().String()}
							if !Contains(backwardCallGraph[fun], currentFun) {
								if funcDecls[fun] == true {
									backwardCallGraph[fun] = append(backwardCallGraph[fun], currentFun)
								}
							}
						}
					}
				case *ast.FuncDecl:
					currentFun = FuncDescriptor{pkg.TypesInfo.Defs[x.Name].Id(), pkg.TypesInfo.Defs[x.Name].Type().String()}
					fmt.Println("\t\t\tFuncDecl:", pkg.TypesInfo.Defs[x.Name].Id(), pkg.TypesInfo.Defs[x.Name].Type().String())
				}
				return true
			})
		}
	}
	return backwardCallGraph
}

func FindFuncDecls(projectPath string, packagePattern string) map[FuncDescriptor]bool {
	fset := token.NewFileSet()
	cfg := &packages.Config{Fset: fset, Mode: mode, Dir: projectPath}
	pkgs, err := packages.Load(cfg, packagePattern)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("FindFuncDecls")
	funcDecls := make(map[FuncDescriptor]bool)
	for _, pkg := range pkgs {
		fmt.Println("\t", pkg)
		for _, node := range pkg.Syntax {
			fmt.Println("\t\t", fset.File(node.Pos()).Name())
			ast.Inspect(node, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.FuncDecl:
					fmt.Println("\t\t\tFuncDecl:", pkg.TypesInfo.Defs[x.Name].Id(), pkg.TypesInfo.Defs[x.Name].Type().String())
					funcDecls[FuncDescriptor{pkg.TypesInfo.Defs[x.Name].Id(), pkg.TypesInfo.Defs[x.Name].Type().String()}] = true
				}
				return true
			})
		}
	}
	return funcDecls
}

func InferRootFunctionsFromGraph(callgraph map[FuncDescriptor][]FuncDescriptor) []FuncDescriptor {
	var allFunctions map[FuncDescriptor]bool
	var rootFunctions []FuncDescriptor
	allFunctions = make(map[FuncDescriptor]bool)
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

func Generatecfg(callgraph map[FuncDescriptor][]FuncDescriptor, path string) {
	functions := make(map[FuncDescriptor]bool, 0)
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
		out.WriteString(f.TypeHash())
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

			out.WriteString(childFun.TypeHash())

			out.WriteString("', ")
			out.WriteString("target: '")
			out.WriteString(k.TypeHash())
			out.WriteString("' ")
			out.WriteString("} },")
			edgeCounter++
		}
	}
	out.WriteString("\n\t]")
	out.WriteString("\n};")
}
