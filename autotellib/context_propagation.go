package autotellib

import (
	"go/ast"
	"go/parser"
	"go/token"

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
		return true
	})
}
