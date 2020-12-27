package frontend

import (
	"net/http"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

type LoginFormData struct {
	Email               string `schema:"email"`
	Password            string `schema:"password"`
	ClientId            string `schema:"client_id"`
	State               string `schema:"state"`
	Scope               string `schema:"scope"`
	RedirectUri         string `schema:"redirect_uri"`
	CodeChallengeMethod string `schema:"code_challenge_method"`
	CodeChallenge       string `schema:"code_challenge"`
	Remember            bool   `schema:"remember"`
	Remember_me         bool   `schema:"remember_me"`
}

func ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/login/index.html")
}
