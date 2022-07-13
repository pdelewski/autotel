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
			var out *os.File
			fmt.Println("\t\t", fset.File(node.Pos()).Name())
			if len(passFileSuffix) > 0 {
				out, _ = os.Create(fset.File(node.Pos()).Name() + passFileSuffix)
				defer out.Close()
			} else {
				out, _ = os.Create(fset.File(node.Pos()).Name() + "ir_context")
				defer out.Close()
			}

			if len(rootFunctions) == 0 {
				printer.Fprint(out, fset, node)
				continue
			}
			astutil.AddImport(fset, node, "context")

			emitCallExpr := func(ident *ast.Ident, n ast.Node, ctxArg *ast.Ident) {
				switch x := n.(type) {
				case *ast.CallExpr:
					if pkg.TypesInfo.Uses[ident] != nil {
						fun := FuncDescriptor{pkg.TypesInfo.Uses[ident].Id(),
							pkg.TypesInfo.Uses[ident].Type().String()}
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
					fun := FuncDescriptor{pkg.TypesInfo.Defs[x.Name].Id(),
						pkg.TypesInfo.Defs[x.Name].Type().String()}

					for k, v := range callgraph {
						if k.TypeHash() == fun.TypeHash() {
							exists = true
						}
						for _, e := range v {
							if fun.TypeHash() == e.TypeHash() {
								exists = true
							}
						}
					}
					if !exists {
						break
					}

					if Contains(rootFunctions, fun) {
						break
					}
					visited := map[FuncDescriptor]bool{}
					fmt.Println("\t\t\tFuncDecl:", pkg.TypesInfo.Defs[x.Name].Id(), pkg.TypesInfo.Defs[x.Name].Type().String())
					if isPath(callgraph, fun, rootFunctions[0], visited) {
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
							fun := FuncDescriptor{pkg.TypesInfo.Defs[method.Names[0]].Id(),
								pkg.TypesInfo.Defs[method.Names[0]].Type().String()}
							fmt.Println("\t\t\tInterfaceType", fun.Id, fun.DeclType)
							if isPath(callgraph, fun, rootFunctions[0], visited) {
								funcType.Params.List = append(funcType.Params.List, ctxField)
							}

						}
					}
				}
				return true
			})
			printer.Fprint(out, fset, node)
			if len(passFileSuffix) > 0 {
				os.Rename(fset.File(node.Pos()).Name(), fset.File(node.Pos()).Name()+".original")
			} else {
				os.Rename(fset.File(node.Pos()).Name()+"ir_context", fset.File(node.Pos()).Name())
			}
		}
	}
}
