package lib

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

func Instrument(projectPath string,
	packagePattern string,
	file string,
	callgraph map[FuncDescriptor][]FuncDescriptor,
	rootFunctions []FuncDescriptor,
	passFileSuffix string) {

	fset := token.NewFileSet()
	fmt.Println("Instrumentation")
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
			astutil.AddNamedImport(fset, node, "otel", "go.opentelemetry.io/otel")

			childTracingTodo := &ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.Ident{
						Name: "__child_tracing_ctx",
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
								Name: "TODO",
							},
						},
						Lparen:   62,
						Ellipsis: 0,
					},
				},
			}
			childTracingSupress := &ast.AssignStmt{
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

			ast.Inspect(node, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.FuncDecl:
					fun := FuncDescriptor{pkg.TypesInfo.Defs[x.Name].Id(), pkg.TypesInfo.Defs[x.Name].Type().String()}
					// that's kind a trick
					// context propagation pass adds additional context parameter
					// this additional parameter has to be removed to match
					// what's already in function callgraph
					fun.DeclType = strings.ReplaceAll(fun.DeclType, "(__tracing_ctx context.Context", "(")
					fun.DeclType = strings.ReplaceAll(fun.DeclType, ", __tracing_ctx context.Context", "")

					// check if it's root function or
					// one of function in call graph
					// and emit proper ast nodes
					_, exists := callgraph[fun]
					if !exists {
						if !Contains(rootFunctions, fun) {
							x.Body.List = append([]ast.Stmt{childTracingTodo, childTracingSupress}, x.Body.List...)
							return false
						}
					}

					for _, root := range rootFunctions {
						visited := map[FuncDescriptor]bool{}

						fmt.Println("\t\t\tFuncDecl:", pkg.TypesInfo.Defs[x.Name].Id(), pkg.TypesInfo.Defs[x.Name].Type().String())
						if isPath(callgraph, fun, root, visited) && fun.TypeHash() != root.TypeHash() {
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
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "span",
										},
										Sel: &ast.Ident{
											Name: "End",
										},
									},
									Lparen:   41,
									Ellipsis: 0,
								},
							}
							_ = s1
							x.Body.List = append([]ast.Stmt{s2, s3, s4}, x.Body.List...)
						} else {
							// check whether this function is root function
							if !Contains(rootFunctions, fun) {
								x.Body.List = append([]ast.Stmt{childTracingTodo, childTracingSupress}, x.Body.List...)
								return false
							}
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
							s3 := &ast.DeferStmt{
								Defer: 27,
								Call: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "rtlib",
										},
										Sel: &ast.Ident{
											Name: "Shutdown",
										},
									},
									Lparen: 48,
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ts",
										},
									},
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

							s8 := &ast.DeferStmt{
								Defer: 27,
								Call: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "span",
										},
										Sel: &ast.Ident{
											Name: "End",
										},
									},
									Lparen:   41,
									Ellipsis: 0,
								},
							}
							_ = s1
							x.Body.List = append([]ast.Stmt{s2, s3, s4, s5, s6, s8}, x.Body.List...)
							x.Body.List = append([]ast.Stmt{childTracingTodo, childTracingSupress}, x.Body.List...)

						}
					}
				case *ast.FuncLit:
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
												Name: `"` + "anonymous" + `"`,
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
										Name: `"` + "anonymous" + `"`,
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
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "span",
								},
								Sel: &ast.Ident{
									Name: "End",
								},
							},
							Lparen:   41,
							Ellipsis: 0,
						},
					}
					_ = s1
					x.Body.List = append([]ast.Stmt{s2, s3, s4}, x.Body.List...)
				}
				return true
			})
			printer.Fprint(out, fset, node)
			os.Rename(fset.File(node.Pos()).Name(), fset.File(node.Pos()).Name()+".tmp")
		}
	}
}
