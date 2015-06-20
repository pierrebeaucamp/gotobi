package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/pierrebeaucamp/gotobi/controllers"
)

var (
	Address = flag.String("address", "", "the address to host on")
	Port    = flag.Int("port", 8000, "the port to host on")
)

func init() {
	var err error
	controllers.DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		os.Exit(1)
	}

	_, err = controllers.DB.Exec(`CREATE TABLE IF NOT EXISTS projects (
		id VARCHAR(100) UNIQUE PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		email VARCHAR(50) UNIQUE NOT NULL,
    	bio TEXT);`)
	if err != nil {
		fmt.Printf("Error: %v \n", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	endpoint := fmt.Sprintf("%v:%v", *Address, *Port)

	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/login/", controllers.PaypalLogin)
	http.HandleFunc("/mail/", controllers.InboundMail)
	http.HandleFunc("/submit/", controllers.Submit)
	http.Handle("/assets/", http.StripPrefix("/assets",
		http.FileServer(http.Dir("assets"))))

	err := http.ListenAndServe(endpoint, nil)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}
