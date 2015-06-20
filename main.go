package main

import (
	"controllers"
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

	http.HandleFunc("/", controllers.Index)
	http.Handle("/assets/", http.StripPrefix("/assets",
		http.FileServer(http.Dir("assets"))))

	err := http.ListenAndServe(endpoint, nil)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
