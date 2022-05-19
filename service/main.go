package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

func main() {

	http.HandleFunc("/inject", inject)

	http.ListenAndServe(":8090", nil)
}
