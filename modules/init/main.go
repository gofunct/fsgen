package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "hello there, "+ os.Getenv("USER"))
	})
	log.Println("starting server on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

