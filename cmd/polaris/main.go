package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/guillaumebchd/polaris/internal/frontend"
	"github.com/guillaumebchd/polaris/internal/oauth"
	"github.com/guillaumebchd/polaris/pkg/token"
)

func main() {

	r := mux.NewRouter()

	/* OAUTH */
	r.HandleFunc("/authorize", oauth.AuthorizeHandler).Methods("GET")
	r.HandleFunc("/token", NotImplemented).Methods("POST")
	r.HandleFunc("/key", token.ServePubKeyHandler).Methods("GET")

	/* USER */
	r.HandleFunc("/login", frontend.ServeLoginPage).Methods("GET")
	r.HandleFunc("/login", oauth.LoginFormHandler).Methods("POST")
	//r.HandleFunc("/login/oauth/authorize", frontend.ServeAuthorizationPage).Methods("GET")
	//r.HandleFunc("/login/oauth/authorize", frontend.AuthorizationPageFormHandler).Methods("POST")

	r.HandleFunc("/error", frontend.ErrorPageHandler).Methods("GET")

	r.HandleFunc("/register", NotImplemented).Methods("GET")
	r.HandleFunc("/register", NotImplemented).Methods("POST")

	r.HandleFunc("/recover", NotImplemented).Methods("GET")
	r.HandleFunc("/recover", NotImplemented).Methods("POST")

	/* CLIENTS */
	r.HandleFunc("/client", NotImplemented).Methods("GET")
	r.HandleFunc("/client", NotImplemented).Methods("POST")

	/* Static */
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
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
