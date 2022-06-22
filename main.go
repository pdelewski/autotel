package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	alib "sumologic.com/autotellib"
)

func usage() {
	fmt.Println("\nusage autotel --command [path to go project] [package pattern]")
	fmt.Println("\tcommand:")
	fmt.Println("\t\tinject                                 (injects open telemetry calls into project code)")
	fmt.Println("\t\tinject-using-graph graph-file          (injects open telemetry calls into project code using provided graph information)")
	fmt.Println("\t\tdumpcfg                                (dumps control flow graph)")
	fmt.Println("\t\tgencfg                                 (generates json representation of control flow graph)")
	fmt.Println("\t\trootfunctions                          (dumps root functions)")
	fmt.Println("\t\trevert                                 (delete generated files)")
}

func inject(root string, packagePattern string) {
	var rootFunctions []string

	rootFunctions = append(rootFunctions, alib.GlobalFindRootFunctions(root, packagePattern)...)

	funcDecls := alib.GlobalFindFuncDecls(root, packagePattern)
	backwardCallGraph := alib.GlobalBuildCallGraph(root, packagePattern, funcDecls)

	alib.ExecutePasses(root, packagePattern, rootFunctions, funcDecls, backwardCallGraph)
}

// Parsing algorithm works as follows. It goes through all function
// decls and infer function bodies to find call to SumoAutoInstrument
// A parent function of this call will become root of instrumentation
// Each function call from this place will be instrumented automatically

func main() {
	fmt.Println("autotel compiler")
	args := len(os.Args)
	if args < 4 {
		usage()
		return
	}
	if os.Args[1] == "--inject" {
		projectPath := os.Args[2]
		packagePattern := os.Args[3]
		inject(projectPath, packagePattern)
		fmt.Println("\tinstrumentation done")
	}
	if os.Args[1] == "--inject-using-graph" {
		graphFile := os.Args[2]
		file, err := os.Open(graphFile)
		if err != nil {
			usage()
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		backwardCallGraph := make(map[string][]string)

		for scanner.Scan() {
			line := scanner.Text()
			keyValue := strings.Split(line, " ")
			funList := []string{}
			fmt.Print("\n\t", keyValue[0])
			for i := 1; i < len(keyValue); i++ {
				fmt.Print(" ", keyValue[i])
				funList = append(funList, keyValue[i])
			}
			backwardCallGraph[keyValue[0]] = funList
		}
		rootFunctions := alib.InferRootFunctionsFromGraph(backwardCallGraph)
		for _, v := range rootFunctions {
			fmt.Println("\nroot:" + v)
		}
		projectPath := os.Args[3]
		packagePattern := os.Args[4]

		funcDecls := alib.GlobalFindFuncDecls(projectPath, packagePattern)

		alib.ExecutePasses(projectPath, packagePattern, rootFunctions, funcDecls, backwardCallGraph)
	}
	if os.Args[1] == "--dumpcfg" {
		projectPath := os.Args[2]
		packagePattern := os.Args[3]
		funcDecls := alib.GlobalFindFuncDecls(projectPath, packagePattern)
		backwardCallGraph := alib.GlobalBuildCallGraph(projectPath, packagePattern, funcDecls)

		fmt.Println("\n\tchild parent")
		for k, v := range backwardCallGraph {
			fmt.Print("\n\t", k)
			fmt.Print(" ", v)
		}
	}
	if os.Args[1] == "--gencfg" {
		projectPath := os.Args[2]
		packagePattern := os.Args[3]
		funcDecls := alib.GlobalFindFuncDecls(projectPath, packagePattern)
		backwardCallGraph := alib.GlobalBuildCallGraph(projectPath, packagePattern, funcDecls)
		alib.Generatecfg(backwardCallGraph, "callgraph.js")
	}
	if os.Args[1] == "--rootfunctions" {
		projectPath := os.Args[2]
		packagePattern := os.Args[3]
		var rootFunctions []string
		rootFunctions = append(rootFunctions, alib.GlobalFindRootFunctions(projectPath, packagePattern)...)
		fmt.Println("rootfunctions:")
		for _, fun := range rootFunctions {
			fmt.Println("\t" + fun)
		}
	}
	if os.Args[1] == "--revert" {
		projectPath := os.Args[2]
		alib.Revert(projectPath)
	}

}
