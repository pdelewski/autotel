package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	alib "sumologic.com/autotellib"
)

func usage() {
	fmt.Println("\nusage autotel --command [path to go project]")
	fmt.Println("\tcommand:")
	fmt.Println("\t\tinject                                 (injects open telemetry calls into project code)")
	fmt.Println("\t\tinject-using-graph graph-file          (injects open telemetry calls into project code using provided graph information)")
	fmt.Println("\t\tdumpcfg                                (dumps control flow graph)")
	fmt.Println("\t\tgencfg                                 (generates json representation of control flow graph)")
	fmt.Println("\t\trootfunctions                          (dumps root functions)")
}

func inject(root string) {
	files := alib.SearchFiles(root)

	var rootFunctions []string
	backwardCallGraph := make(map[string]string)

	for _, file := range files {
		rootFunctions = append(rootFunctions, alib.FindRootFunctions(file)...)
	}
	for _, file := range files {
		callGraphInstance := alib.BuildCallGraph(file)
		for key, value := range callGraphInstance {
			backwardCallGraph[key] = value
		}
	}
	for _, file := range files {
		alib.Instrument(file, backwardCallGraph, rootFunctions)
	}
}

// Parsing algorithm works as follows. It goes through all function
// decls and infer function bodies to find call to SumoAutoInstrument
// A parent function of this call will become root of instrumentation
// Each function call from this place will be instrumented automatically

func main() {
	fmt.Println("autotel compiler")
	args := len(os.Args)
	if args < 3 {
		usage()
		return
	}
	if os.Args[1] == "--inject" {
		inject(os.Args[2])
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
		backwardCallGraph := make(map[string]string)
		for scanner.Scan() {
			line := scanner.Text()
			keyValue := strings.Split(line, " ")
			fmt.Print("\n\t", keyValue[0])
			fmt.Print(" ", keyValue[1])

			backwardCallGraph[keyValue[0]] = keyValue[1]
		}
		rootFunctions := alib.InferRootFunctionsFromGraph(backwardCallGraph)
		for _, v := range rootFunctions {
			fmt.Println("\nroot:" + v)
		}
		files := alib.SearchFiles(os.Args[3])
		for _, file := range files {
			alib.Instrument(file, backwardCallGraph, rootFunctions)
		}
	}
	if os.Args[1] == "--dumpcfg" {
		files := alib.SearchFiles(os.Args[2])
		backwardCallGraph := make(map[string]string)
		for _, file := range files {
			callGraphInstance := alib.BuildCallGraph(file)
			for key, value := range callGraphInstance {
				backwardCallGraph[key] = value
			}
		}
		fmt.Println("\n\tchild parent")
		for k, v := range backwardCallGraph {
			fmt.Print("\n\t", k)
			fmt.Print(" ", v)
		}
	}
	if os.Args[1] == "--gencfg" {
		files := alib.SearchFiles(os.Args[2])
		backwardCallGraph := make(map[string]string)
		for _, file := range files {
			callGraphInstance := alib.BuildCallGraph(file)
			for key, value := range callGraphInstance {
				backwardCallGraph[key] = value
			}
		}
		alib.Generatecfg(backwardCallGraph)
	}
	if os.Args[1] == "--rootfunctions" {
		files := alib.SearchFiles(os.Args[2])
		var rootFunctions []string
		for _, file := range files {
			rootFunctions = append(rootFunctions, alib.FindRootFunctions(file)...)
		}
		for _, fun := range rootFunctions {
			fmt.Println("\t" + fun)
		}
	}

}
