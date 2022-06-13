package autotellib

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func PropagateContext(file string, callgraph map[string][]string, rootFunctions []string, funcDecls map[string]bool, passFileSuffix string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}
	out, _ := os.Create(file + passFileSuffix)
	defer out.Close()

	if len(rootFunctions) == 0 {
		printer.Fprint(out, fset, node)
		return
	}

	astutil.AddImport(fset, node, "context")

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
			for k, v := range callgraph {
				if k == x.Name.Name {
					exists = true
				}
				for _, e := range v {
					if x.Name.Name == e {
						exists = true
					}
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
					x.Args = append(x.Args, ctxArg)
				}
			}
			_, ok = x.Fun.(*ast.FuncLit)
			if ok {
				x.Args = append(x.Args, ctxArg)
			}
			// TODO selectors are recursive
			// to handle a.b.c.fun()
			// all selectors have to unpacked
			sel, ok := x.Fun.(*ast.SelectorExpr)
			if ok {
				// packageIdent, ok := sel.X.(*ast.Ident)
				found := funcDecls[sel.Sel.Name]
				if found {
					x.Args = append(x.Args, ctxArg)
				}
			}
		case *ast.FuncLit:
			x.Type.Params.List = append(x.Type.Params.List, ctxField)
		}
		return true
	})

	printer.Fprint(out, fset, node)
}
