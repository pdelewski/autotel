package autotellib

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func PropagateContext(file string, callgraph map[string]string, rootFunctions []string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	astutil.AddImport(fset, node, "context")
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// inject context only
			// functions available in the call graph
			_, exists := callgraph[x.Name.Name]
			if !exists {
				return false
			}
			fmt.Printf("function decl: %s, parameters:\n", x.Name)
			for _, param := range x.Type.Params.List {
				fmt.Printf("  Name: %s\n", param.Names[0])
				fmt.Printf("    ast type          : %T\n", param.Type)
				fmt.Printf("    type desc         : %+v\n", param.Type)
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

			x.Type.Params.List = append([]*ast.Field{ctxField}, x.Type.Params.List...)
		case *ast.CallExpr:
			ident, ok := x.Fun.(*ast.Ident)
			if ok {
				fmt.Println("call:", ident.Name)
			}
		}
		return true
	})
	out, err := os.Create(file + ".pass_ctx")
	defer out.Close()

	printer.Fprint(out, fset, node)
}
