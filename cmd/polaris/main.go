package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/guillaumebchd/polaris/internal/frontend"
	"github.com/guillaumebchd/polaris/internal/oauth"
	"github.com/guillaumebchd/polaris/pkg/code"
	"github.com/guillaumebchd/polaris/pkg/reset"
	"github.com/guillaumebchd/polaris/pkg/token"
	"github.com/joho/godotenv"
)

func main() {

	/*
		Setup of pkg
	*/
	godotenv.Load()
	code.Initialize()
	token.Initialize()
	reset.Initialize()

	r := mux.NewRouter()

	r.HandleFunc("/", frontend.ServeHomePage).Methods("GET")

	/* OAUTH */
	r.HandleFunc("/authorize", oauth.AuthorizeHandler).Methods("GET")
	r.HandleFunc("/token", oauth.TokenHandler).Methods("POST")
	r.HandleFunc("/key", token.ServePubKeyHandler).Methods("GET")

	/* USER */
	r.HandleFunc("/login", frontend.ServeLoginPage).Methods("GET")
	r.HandleFunc("/login", oauth.LoginFormHandler).Methods("POST")

	r.HandleFunc("/error", frontend.ErrorPageHandler).Methods("GET")

	r.HandleFunc("/register", frontend.ServeRegisterPage).Methods("GET")
	r.HandleFunc("/register", frontend.RegisterFormHandler).Methods("POST")

	r.HandleFunc("/recover", frontend.ServeRecoverPage).Methods("GET")
	r.HandleFunc("/recover", frontend.RecoverFormHandler).Methods("POST")

	r.HandleFunc("/reset/{code}", frontend.ServeResetPage).Methods("GET")
	r.HandleFunc("/reset/{code}", frontend.ServeRecoverPage).Methods("POST")

	/* CLIENTS */
	r.HandleFunc("/client", NotImplemented).Methods("GET")
	r.HandleFunc("/client", NotImplemented).Methods("POST")

	/* Static */
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Starting server on : " + srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func Test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{ \"oui\":\"oui\" }"))
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
}
