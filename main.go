package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	Address = flag.String("address", "", "the address to host on")
	Port    = flag.Int("port", 8000, "the port to host on")
)

func main() {
	flag.Parse()
	endpoint := fmt.Sprintf("%v:%v", *Address, *Port)

	http.HandleFunc("/", index)

	err := http.ListenAndServe(endpoint, nil)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}
