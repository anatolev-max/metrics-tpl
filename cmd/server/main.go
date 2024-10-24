package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/get/{type}`, GetHandler)
	mux.HandleFunc(`/update/{type}/{name}/{value}`, PostHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
