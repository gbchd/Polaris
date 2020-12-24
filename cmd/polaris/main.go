package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/authorize", NotImplemented).Methods("GET")
	r.HandleFunc("/token", NotImplemented).Methods("POST")

	r.HandleFunc("/login", NotImplemented).Methods("GET")
	r.HandleFunc("/login", NotImplemented).Methods("POST")

	r.HandleFunc("/register", NotImplemented).Methods("GET")
	r.HandleFunc("/register", NotImplemented).Methods("POST")

	r.HandleFunc("/client", NotImplemented).Methods("GET")
	r.HandleFunc("/client", NotImplemented).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Starting server on : " + srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
}
