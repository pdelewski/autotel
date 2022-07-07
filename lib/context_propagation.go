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
	callgraph map[FuncDescriptor][]FuncDescriptor,
	rootFunctions []FuncDescriptor,
	funcDecls map[FuncDescriptor]bool,
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

			emitCallExpr := func(ident *ast.Ident, n ast.Node, ctxArg *ast.Ident) {
				switch x := n.(type) {
				case *ast.CallExpr:
					fun := FuncDescriptor{pkg.TypesInfo.Uses[ident].Id(), pkg.TypesInfo.Uses[ident].Type().String()}
					found := funcDecls[fun]
					// inject context parameter only
					// to these functions for which function decl
					// exists

					if found {
						visited := map[FuncDescriptor]bool{}
						if isPath(callgraph, fun, rootFunctions[0], visited) {
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
					funName := FuncDescriptor{pkg.TypesInfo.Defs[x.Name].Id(), pkg.TypesInfo.Defs[x.Name].Type().String()}

					for k, v := range callgraph {
						if k.TypeHash() == funName.TypeHash() {
							exists = true
						}
						for _, e := range v {
							if funName.TypeHash() == e.TypeHash() {
								exists = true
							}
						}
					}
					if !exists {
						break
					}

					if Contains(rootFunctions, FuncDescriptor{pkg.TypesInfo.Defs[x.Name].Id(), pkg.TypesInfo.Defs[x.Name].Type().String()}) {
						break
					}
					visited := map[FuncDescriptor]bool{}
					fmt.Println("\t\t\tFuncDecl:", pkg.TypesInfo.Defs[x.Name].Id(), pkg.TypesInfo.Defs[x.Name].Type().String())
					if isPath(callgraph, FuncDescriptor{pkg.TypesInfo.Defs[x.Name].Id(), pkg.TypesInfo.Defs[x.Name].Type().String()}, rootFunctions[0], visited) {
						x.Type.Params.List = append(x.Type.Params.List, ctxField)
					}
				case *ast.CallExpr:
					ident, ok := x.Fun.(*ast.Ident)

					if ok {
						fmt.Println("\t\t\tCallExpr:", pkg.TypesInfo.Uses[ident].Id(), pkg.TypesInfo.Uses[ident].Type().String())

						emitCallExpr(ident, n, ctxArg)
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
						emitCallExpr(sel.Sel, n, ctxArg)
					}
				case *ast.FuncLit:
					x.Type.Params.List = append(x.Type.Params.List, ctxField)
				case *ast.InterfaceType:
					for _, method := range x.Methods.List {
						if funcType, ok := method.Type.(*ast.FuncType); ok {
							visited := map[FuncDescriptor]bool{}
							fmt.Println("\t\t\tInterfaceType", pkg.TypesInfo.Defs[method.Names[0]].Id(), pkg.TypesInfo.Defs[method.Names[0]].Type().String())
							if isPath(callgraph, FuncDescriptor{pkg.TypesInfo.Defs[method.Names[0]].Id(), pkg.TypesInfo.Defs[method.Names[0]].Type().String()}, rootFunctions[0], visited) {
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
