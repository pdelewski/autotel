package autotellib

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

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

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func Instrument(file string, callgraph map[string]string, rootFunctions []string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}
	astutil.AddImport(fset, node, "context")
	astutil.AddNamedImport(fset, node, "otel", "go.opentelemetry.io/otel")
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// check if it's root function or
			// one of function in call graph
			// and emit proper ast nodes
			_, exists := callgraph[x.Name.Name]
			if !exists {
				if !Contains(rootFunctions, x.Name.Name) {
					return false
				}
			}

			for _, root := range rootFunctions {
				if isPath(callgraph, x.Name.Name, root) && x.Name.Name != root {
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
									Value: `"child instrumentation"`,
								},
							},
						},
					}
					s2 := &ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "__child_tracing_ctx",
							},
							&ast.Ident{
								Name: "span",
							},
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "otel",
											},
											Sel: &ast.Ident{
												Name: "Tracer",
											},
										},
										Lparen: 50,
										Args: []ast.Expr{
											&ast.Ident{
												Name: `"` + x.Name.Name + `"`,
											},
										},
										Ellipsis: 0,
									},
									Sel: &ast.Ident{
										Name: "Start",
									},
								},
								Lparen: 62,
								Args: []ast.Expr{
									&ast.Ident{
										Name: "__tracing_ctx",
									},
									&ast.Ident{
										Name: `"` + x.Name.Name + `"`,
									},
								},
								Ellipsis: 0,
							},
						},
					}

					s3 := &ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "_",
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.Ident{
								Name: "__child_tracing_ctx",
							},
						},
					}

					s4 := &ast.DeferStmt{
						Defer: 27,
						Call: &ast.CallExpr{
							Fun: &ast.FuncLit{
								Type: &ast.FuncType{
									Func:   33,
									Params: &ast.FieldList{},
								},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										&ast.ExprStmt{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "span",
													},
													Sel: &ast.Ident{
														Name: "End",
													},
												},
												Lparen:   49,
												Ellipsis: 0,
											},
										},
									},
								},
							},
							Lparen:   52,
							Ellipsis: 0,
						},
					}
					x.Body.List = append([]ast.Stmt{s1, s2, s3, s4}, x.Body.List...)
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
					s4 := &ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "otel",
								},
								Sel: &ast.Ident{
									Name: "SetTracerProvider",
								},
							},
							Lparen: 49,
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "ts",
									},
									Sel: &ast.Ident{
										Name: "Tp",
									},
								},
							},
							Ellipsis: 0,
						},
					}
					s5 := &ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "ctx",
							},
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "context",
									},
									Sel: &ast.Ident{
										Name: "Background",
									},
								},
								Lparen:   52,
								Ellipsis: 0,
							},
						},
					}
					s6 := &ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "__child_tracing_ctx",
							},
							&ast.Ident{
								Name: "span",
							},
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "otel",
											},
											Sel: &ast.Ident{
												Name: "Tracer",
											},
										},
										Lparen: 50,
										Args: []ast.Expr{
											&ast.Ident{
												Name: `"` + x.Name.Name + `"`,
											},
										},
										Ellipsis: 0,
									},
									Sel: &ast.Ident{
										Name: "Start",
									},
								},
								Lparen: 62,
								Args: []ast.Expr{
									&ast.Ident{
										Name: "ctx",
									},
									&ast.Ident{
										Name: `"` + x.Name.Name + `"`,
									},
								},
								Ellipsis: 0,
							},
						},
					}
					s7 := &ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.Ident{
								Name: "_",
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.Ident{
								Name: "__child_tracing_ctx",
							},
						},
					}
					s8 := &ast.DeferStmt{
						Defer: 27,
						Call: &ast.CallExpr{
							Fun: &ast.FuncLit{
								Type: &ast.FuncType{
									Func:   33,
									Params: &ast.FieldList{},
								},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										&ast.ExprStmt{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "span",
													},
													Sel: &ast.Ident{
														Name: "End",
													},
												},
												Lparen:   49,
												Ellipsis: 0,
											},
										},
									},
								},
							},
							Lparen:   52,
							Ellipsis: 0,
						},
					}
					x.Body.List = append([]ast.Stmt{s1, s2, s3, s4, s5, s6, s7, s8}, x.Body.List...)
				}
			}

		}
		return true
	})

	out, err := os.Create(file + ".pass_tracing")
	defer out.Close()

	printer.Fprint(out, fset, node)
}
