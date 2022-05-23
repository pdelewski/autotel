package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	alib "sumologic.com/autotellib"
)

var projectDir string

func readGraphBody(graphFile string) {
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
	if len(rootFunctions) != 1 {
		panic("more than one graph")
	}
	for _, v := range rootFunctions {
		fmt.Println("\nroot:" + v)
	}
	files := alib.SearchFiles(projectDir)
	for _, file := range files {
		alib.Instrument(file, backwardCallGraph, rootFunctions)
	}
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
	fmt.Println("\nusage autotelservice [path to go project]")
}

func main() {
	args := len(os.Args)
	if args < 2 {
		usage()
		return
	}
	files := alib.SearchFiles(os.Args[1])
	projectDir = os.Args[1]
	backwardCallGraph := make(map[string]string)
	for _, file := range files {
		callGraphInstance := alib.BuildCallGraph(file)
		for key, value := range callGraphInstance {
			backwardCallGraph[key] = value
		}
	}
	alib.Generatecfg(backwardCallGraph, "./static/callgraph.js")

	http.HandleFunc("/inject", inject)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.ListenAndServe(":8090", nil)
}
