package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	alib "github.com/pdelewski/autotel/lib"
)

var projectDir string
var packagePattern string

func readGraphBody(graphFile string) {
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
	if len(rootFunctions) != 1 {
		panic("more than one graph")
	}
	for _, v := range rootFunctions {
		fmt.Println("\nroot:" + v)
	}

	funcDecls := alib.FindFuncDecls(projectDir, packagePattern)

	alib.ExecutePasses(projectDir,
		packagePattern,
		rootFunctions,
		funcDecls,
		backwardCallGraph)
}

func inject(w http.ResponseWriter, r *http.Request) {
	var bodyBytes []byte
	var err error

	if r.Body != nil {
		bodyBytes, err = ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Body reading error: %v", err)
			return
		}
		defer r.Body.Close()
	}

	if len(bodyBytes) > 0 {
		fmt.Println(string(bodyBytes))
		f, err := os.Create("graphbody")

		if err != nil {
			fmt.Println(err)
		}

		defer f.Close()

		_, errSave := f.WriteString(string(bodyBytes))
		if errSave != nil {
			fmt.Println(errSave)
		}
		readGraphBody("graphBody")
	} else {
		fmt.Printf("Body: No Body Supplied\n")
	}
	fmt.Fprintf(w, "inject\n")

}

func usage() {
	fmt.Println("\nusage autotelservice [path to go project] [package pattern]")
}

func main() {
	args := len(os.Args)
	if args < 3 {
		usage()
		return
	}

	projectDir = os.Args[1]
	packagePattern = os.Args[2]

	funcDecls := alib.FindFuncDecls(projectDir, packagePattern)

	backwardCallGraph := alib.BuildCallGraph(projectDir,
		packagePattern,
		funcDecls)
	alib.Generatecfg(backwardCallGraph, "./static/callgraph.js")

	http.HandleFunc("/inject", inject)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.ListenAndServe(":8090", nil)
}
