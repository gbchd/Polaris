package frontend

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guillaumebchd/polaris/pkg/reset"
)

type ResetFormData struct {
	Password  string `schema:"password"`
	VPassword string `schema:"vpassword"`
}

func ServeResetPage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	code := params["code"]
	fmt.Println(code)

	if reset.Exist(code) {
		http.ServeFile(w, r, "web/reset/index.html")
	} else {
		w.Write([]byte("The code does not exist"))
	}

}

func ResetFormHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	code := params["code"]
	fmt.Println(code)

	if reset.Exist(code) {
		w.Write([]byte("The code exists"))
	} else {
		w.Write([]byte("The code does not exist"))
	}
}
