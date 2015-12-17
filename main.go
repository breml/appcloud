package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Use port: %s\n", port)

	http.HandleFunc("/", helloHandler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Panicf("ListenAndServe error: %v\n", err)
	}
}
