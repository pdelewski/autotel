package autotellib

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func PropagateContext(file string, callgraph map[string]string, rootFunctions []string, funcDecls map[string]bool) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}
	out, _ := os.Create(file + ".pass_ctx")
	defer out.Close()

	if len(rootFunctions) == 0 {
		printer.Fprint(out, fset, node)
		return
	}

	astutil.AddImport(fset, node, "context")

	ast.Inspect(node, func(n ast.Node) bool {
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
			for k, v := range callgraph {
				if k == x.Name.Name || v == x.Name.Name {
					exists = true
				}
			}
			if !exists {
				break
			}

			if Contains(rootFunctions, x.Name.Name) {
				break
			}
			// fmt.Printf("function decl: %s, parameters:\n", x.Name)
			// for _, param := range x.Type.Params.List {
			// fmt.Printf("  Name: %s\n", param.Names[0])
			// fmt.Printf("    ast type          : %T\n", param.Type)
			// fmt.Printf("    type desc         : %+v\n", param.Type)
			// }
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
			x.Type.Params.List = append(x.Type.Params.List, ctxField)
		case *ast.CallExpr:
			ident, ok := x.Fun.(*ast.Ident)
			if ok {
				found := funcDecls[ident.Name]
				_ = found
				// inject context parameter only
				// to these functions for which function decl
				// exists

				if found {
					ctxArg := &ast.Ident{
						Name: "__child_tracing_ctx",
					}
					x.Args = append(x.Args, ctxArg)
				}
			}
			_, ok = x.Fun.(*ast.FuncLit)
			if ok {
				ctxArg := &ast.Ident{
					Name: "__child_tracing_ctx",
				}
				x.Args = append(x.Args, ctxArg)
			}
		case *ast.FuncLit:
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

			x.Type.Params.List = append(x.Type.Params.List, ctxField)
		}
		return true
	})

	printer.Fprint(out, fset, node)
}
