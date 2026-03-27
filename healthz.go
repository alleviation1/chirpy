package main

import (
	"net/http"
)

func Healthz (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	data := []byte("200 OK")
	w.Write(data)
}