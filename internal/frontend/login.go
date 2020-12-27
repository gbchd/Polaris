package frontend

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
	auth "github.com/guillaumebchd/polaris/pkg/authentication"
	"github.com/guillaumebchd/polaris/pkg/code"
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

func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var data LoginFormData
	err = decoder.Decode(&data, r.Form)
	if err != nil {
		log.Fatal("Error in GET parameters : ", err)
	}

	// If the client_id is missing
	if data.ClientId == "" {
		u := data.RedirectUri
		if u == "" {
			u = "/error"
		}

		ur, err := url.Parse(u)
		if err != nil {
			panic(err)
		}

		q := ur.Query()
		q.Add("error", "missing_client_id")
		q.Add("error_description", "The client_id is missing from the request")
		ur.RawQuery = q.Encode()

		http.Redirect(w, r, ur.String(), http.StatusSeeOther)
		return
	}

	client, err := auth.GetClient(data.ClientId)
	if err != nil {
		u := data.RedirectUri
		if u == "" {
			u = "/error"
		}

		ur, err := url.Parse(u)
		if err != nil {
			panic(err)
		}

		q := ur.Query()
		q.Add("error", "invalid_client_id")
		q.Add("error_description", "The client_id given in the request is not valid")
		ur.RawQuery = q.Encode()

		http.Redirect(w, r, ur.String(), http.StatusSeeOther)
		return
	}

	if data.RedirectUri == "" {
		data.RedirectUri = client.RedirectUri
	}

	ur, err := url.Parse(data.RedirectUri)
	if err != nil {
		panic(err)
	}

	q := ur.Query()

	// Check auth
	if !auth.CheckUser(data.Email, data.Password) {
		q.Add("error", "access_denied")
		q.Add("error_description", "The user did not consent.")
		ur.RawQuery = q.Encode()

		http.Redirect(w, r, ur.String(), http.StatusSeeOther)
		return
	}

	c_data := code.CodeData{
		Openid:     strings.Contains(data.Scope, "openid"),
		PKCEMethod: data.CodeChallengeMethod,
		PKCEValue:  data.CodeChallenge,
	}

	c, err := code.GenerateCode(c_data)
	if err != nil {
		q.Add("error", "internal_error")
		q.Add("error_description", "Something went wrong when creating your code.")
		ur.RawQuery = q.Encode()

		http.Redirect(w, r, ur.String(), http.StatusSeeOther)
		return
	}

	q.Add("code", c)
	if data.State != "" {
		q.Add("state", data.State)
	}

	ur.RawQuery = q.Encode()
	http.Redirect(w, r, ur.String(), 302)
}

func ServeAuthorizationPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
}

func AuthorizationPageFormHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
}
