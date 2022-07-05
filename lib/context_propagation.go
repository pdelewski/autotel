package lib

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"os"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

func PropagateContext(projectPath string,
	packagePattern string,
	callgraph map[string][]string,
	rootFunctions []string,
	funcDecls map[string]bool,
	passFileSuffix string) {

	fset := token.NewFileSet()
	fmt.Println("PropagateContext")
	cfg := &packages.Config{Fset: fset, Mode: mode, Dir: projectPath}
	pkgs, err := packages.Load(cfg, packagePattern)
	if err != nil {
		log.Fatal(err)
	}
	for _, pkg := range pkgs {
		fmt.Println("\t", pkg)

		for _, node := range pkg.Syntax {
			fmt.Println("\t\t", fset.File(node.Pos()).Name())
			out, _ := os.Create(fset.File(node.Pos()).Name() + passFileSuffix)
			defer out.Close()

			if len(rootFunctions) == 0 {
				printer.Fprint(out, fset, node)
				continue
			}
			astutil.AddImport(fset, node, "context")

			emitCallExpr := func(name string, n ast.Node, ctxArg *ast.Ident) {
				switch x := n.(type) {
				case *ast.CallExpr:
					found := funcDecls[name]
					// inject context parameter only
					// to these functions for which function decl
					// exists

					if found {
						visited := map[string]bool{}
						if isPath(callgraph, name, rootFunctions[0], visited) {
							x.Args = append(x.Args, ctxArg)
						}
					}
				}
			}

			ast.Inspect(node, func(n ast.Node) bool {
				ctxArg := &ast.Ident{
					Name: "__child_tracing_ctx",
				}
				ctxField := &ast.Field{
					Names: []*ast.Ident{
						&ast.Ident{
							Name: "__tracing_ctx",
						},
					},
					Type: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "context",
						},
						Sel: &ast.Ident{
							Name: "Context",
						},
					},
				}
				switch x := n.(type) {
				case *ast.FuncDecl:
					// inject context only
					// functions available in the call graph
					// _, exists := callgraph[x.Name.Name]
					// if !exists {
					// 	return false
					// }
					// TODO this is not optimap o(n)
					exists := false
					funName := x.Name.Name
					for k, v := range callgraph {
						if k == funName {
							exists = true
						}
						for _, e := range v {
							if funName == e {
								exists = true
							}
						}
					}
					if !exists {
						break
					}

					if Contains(rootFunctions, funName) {
						break
					}
					visited := map[string]bool{}
					fmt.Println("\t\t\tFuncDecl:", pkg.TypesInfo.Defs[x.Name].Id(), pkg.TypesInfo.Defs[x.Name].Type().String())
					if isPath(callgraph, funName, rootFunctions[0], visited) {
						x.Type.Params.List = append(x.Type.Params.List, ctxField)
					}
				case *ast.CallExpr:
					ident, ok := x.Fun.(*ast.Ident)

					if ok {
						fmt.Println("\t\t\tCallExpr:", pkg.TypesInfo.Defs[ident].Id(), pkg.TypesInfo.Defs[ident].Type().String())

						emitCallExpr(ident.Name, n, ctxArg)
					}
					_, ok = x.Fun.(*ast.FuncLit)
					if ok {
						x.Args = append(x.Args, ctxArg)
					}
					// TODO selectors are recursive
					// a.b.c.fun()
					// check whether the most outer one is package
					sel, ok := x.Fun.(*ast.SelectorExpr)

					if ok {
						emitCallExpr(sel.Sel.Name, n, ctxArg)
					}
				case *ast.FuncLit:
					x.Type.Params.List = append(x.Type.Params.List, ctxField)
				case *ast.InterfaceType:
					for _, method := range x.Methods.List {
						if funcType, ok := method.Type.(*ast.FuncType); ok {
							visited := map[string]bool{}
							fmt.Println("\t\t\tInterfaceType", pkg.TypesInfo.Defs[method.Names[0]].Id(), pkg.TypesInfo.Defs[method.Names[0]].Type().String())
							if isPath(callgraph, method.Names[0].Name, rootFunctions[0], visited) {
								funcType.Params.List = append(funcType.Params.List, ctxField)
							}

						}
					}
				}
				return true
			})
			printer.Fprint(out, fset, node)
			os.Rename(fset.File(node.Pos()).Name(), fset.File(node.Pos()).Name()+".original")
		}
	}
}
