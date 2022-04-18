package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/ast/astutil"
)

func usage() {
	fmt.Println("\nusage autotel [path to go project]")
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

func instrument(file string, callgraph map[string]string, rootFunctions []string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	astutil.AddImport(fset, node, "context")
	astutil.AddNamedImport(fset, node, "otel", "go.opentelemetry.io/otel")
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			id, ok := x.Fun.(*ast.Ident)
			if ok {
				fmt.Println(id)
			}
		case *ast.FuncDecl:
			// check if it's root function or
			// one of function in call graph
			// and emit proper ast nodes
			for _, root := range rootFunctions {
				if isPath(callgraph, x.Name.Name, root) && x.Name.Name != root {
					fmt.Printf("\nInstrument child : %s %s\n", x.Name.Name, root)
					newCallStmt := &ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "fmt",
								},
								Sel: &ast.Ident{
									Name: "Println",
								},
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: `"child instrumentation"`,
								},
							},
						},
					}
					x.Body.List = append([]ast.Stmt{newCallStmt}, x.Body.List...)
				} else {
					s1 := &ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "fmt",
								},
								Sel: &ast.Ident{
									Name: "Println",
								},
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: `"root instrumentation"`,
								},
							},
						},
					}

					s2 :=
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "ts",
								},
							},
							Tok: token.DEFINE,

							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "rtlib",
										},
										Sel: &ast.Ident{
											Name: "NewTracingState",
										},
									},
									Lparen:   54,
									Ellipsis: 0,
								},
							},
						}
					s3 :=
						&ast.DeferStmt{
							Defer: 27,
							Call: &ast.CallExpr{
								Fun: &ast.FuncLit{
									Type: &ast.FuncType{
										Func:   33,
										Params: &ast.FieldList{},
									},
									Body: &ast.BlockStmt{
										List: []ast.Stmt{
											&ast.IfStmt{
												If: 41,
												Init: &ast.AssignStmt{
													Lhs: []ast.Expr{
														&ast.Ident{
															Name: "err",
														},
													},
													Tok: token.DEFINE,
													Rhs: []ast.Expr{
														&ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "ts",
																	},
																	Sel: &ast.Ident{
																		Name: "Tp",
																	},
																},
																Sel: &ast.Ident{
																	Name: "Shutdown",
																},
															},
															Lparen: 65,
															Args: []ast.Expr{
																&ast.CallExpr{
																	Fun: &ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "context",
																		},
																		Sel: &ast.Ident{
																			Name: "Backgroud",
																		},
																	},
																	Lparen:   83,
																	Ellipsis: 0,
																},
															},
															Ellipsis: 0,
														},
													},
												},
												Cond: &ast.BinaryExpr{
													X: &ast.Ident{
														Name: "err",
													},
													OpPos: 92,
													Op:    token.NEQ,
													Y: &ast.Ident{
														Name: "nil",
													},
												},
												Body: &ast.BlockStmt{
													List: []ast.Stmt{
														&ast.ExprStmt{
															X: &ast.CallExpr{
																Fun: &ast.SelectorExpr{
																	X: &ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "ts",
																		},
																		Sel: &ast.Ident{
																			Name: "Logger",
																		},
																	},
																	Sel: &ast.Ident{
																		Name: "Fatal",
																	},
																},
																Lparen: 115,
																Args: []ast.Expr{
																	&ast.Ident{
																		Name: "err",
																	},
																},
																Ellipsis: 0,
															},
														},
													},
												},
											},
										},
									},
								},
								Lparen:   122,
								Ellipsis: 0,
							},
						}
					x.Body.List = append([]ast.Stmt{s1, s2, s3}, x.Body.List...)
				}
			}

		}
		return true
	})

	out, err := os.Create(file + ".out")
	defer out.Close()

	fmt.Println("Instrumentation result:")
	printer.Fprint(out, fset, node)
}

func parsePath(root string) {
	fmt.Println("parsing", root)
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
	fmt.Println("Root Functions:")
	for _, fun := range rootFunctions {
		fmt.Println(fun)
	}
	fmt.Println("BackwardCallGraph:")
	for k, v := range backwardCallGraph {
		fmt.Println(k, v)
	}
	fmt.Println("Instrument:")
	for _, file := range files {
		instrument(file, backwardCallGraph, rootFunctions)
	}
}

// Parsing algorithm works as follows. It goes through all function
// decls and infer function bodies to find call to SumoAutoInstrument
// A parent function of this call will become root of instrumentation
// Each function call from this place will be instrumented automatically

func main() {
	fmt.Println("autotel compiler")
	args := len(os.Args)
	if args != 2 {
		usage()
		return
	}
	parsePath(os.Args[1])

}
