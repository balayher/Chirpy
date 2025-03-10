package main

import (
	"log"
	"net/http"
)

func main() {
	ServeMux := http.NewServeMux()
	ServeMux.Handle("/", http.FileServer(http.Dir(".")))
	server := http.Server{
		Handler: ServeMux,
		Addr:    ":8080",
	}

	log.Printf("Serving files from . on port 8080")
	log.Fatal(server.ListenAndServe())
}
