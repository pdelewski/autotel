package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	alib "sumologic.com/autotellib"
)

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
