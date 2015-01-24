package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	PORT = os.Getenv("PORT")
)

func indexHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "hello, world")
}

func main() {
	// index handler
	http.HandleFunc("/", indexHandler)

	// server listener
	log.Printf("Listening on :%s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
