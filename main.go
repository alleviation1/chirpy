package main

import (
	"net/http"
	"log"
)

func main() {
	serverMux := http.NewServeMux()
	serverMux.Handle("/", http.FileServer(http.Dir(".")))

	server := &http.Server{
		Addr: ":8080",
		Handler: serverMux,
	}
	log.Fatal(server.ListenAndServe())

	defer server.Close()
}