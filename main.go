package main

import (
	"net/http"
	"log"
)

func main() {
	const port = ":8080"
	const fileRoot = "."

	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(fileRoot))))

	serverMux.HandleFunc("/healthz", Healthz)

	server := &http.Server{
		Addr: port,
		Handler: serverMux,
	}

	log.Fatal(server.ListenAndServe())

	defer server.Close()
}