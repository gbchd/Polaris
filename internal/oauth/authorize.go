package oauth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
	auth "github.com/guillaumebchd/polaris/pkg/authentication"
)

const (
	ResponseTypeImplicit = "token"
	ResponseTypeCode     = "code"
)

type AuthorizeData struct {
	ResponseType        string `schema:"response_type"`
	ClientId            string `schema:"client_id"`
	Scope               string `schema:"scope"`
	RedirectUri         string `schema:"redirect_uri"`
	State               string `schema:"state"`
	CodeChallengeMethod string `schema:"code_challenge_method"`
	CodeChallenge       string `schema:"code_challenge"`
}

type LoginFormData struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`

	Remember    bool `schema:"remember"`
	Remember_me bool `schema:"remember_me"`

	ResponseType        string `schema:"response_type"`
	ClientId            string `schema:"client_id"`
	Scope               string `schema:"scope"`
	RedirectUri         string `schema:"redirect_uri"`
	State               string `schema:"state"`
	CodeChallengeMethod string `schema:"code_challenge_method"`
	CodeChallenge       string `schema:"code_challenge"`
}

var decoder = schema.NewDecoder()

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	var data AuthorizeData
	err := decoder.Decode(&data, r.URL.Query())
	if err != nil {
		fmt.Println("Error in GET parameters : ", err)
	}

	// We check that the client id is here and is not invalid
	if data.ClientId == "" {
		RedirectToError(w, r, data.RedirectUri, "missing_client_id", "The client_id is missing from the request")
		return
	}

	client, err := auth.GetClient(data.ClientId)
	if err != nil {
		RedirectToError(w, r, data.RedirectUri, "invalid_client_id", "The client_id given in the request is not valid")
		return
	}

	if data.RedirectUri == "" {
		data.RedirectUri = client.RedirectUri
	}

	login, err := url.Parse("/login")
	if err != nil {
		panic(err)
	}

	login.RawQuery = r.URL.RawQuery

	switch data.ResponseType {
	case ResponseTypeCode:
		http.Redirect(w, r, login.String(), http.StatusSeeOther)
		return
	case ResponseTypeImplicit:
		http.Redirect(w, r, login.String(), http.StatusSeeOther)
		return
	default:
		RedirectToError(w, r, data.RedirectUri, "invalid_response_type", "The response_type given in the request is not valid")
		return
	}
}

// TODO: Allow the possibility to retry to connect if wrong password (but add captcha after x try)
func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var data LoginFormData
	err = decoder.Decode(&data, r.Form)
	if err != nil {
		fmt.Println("Error in GET parameters : ", err)
	}

	// We check that the client id is here and is not invalid
	if data.ClientId == "" {
		RedirectToError(w, r, data.RedirectUri, "missing_client_id", "The client_id is missing from the request")
		return
	}
	client, err := auth.GetClient(data.ClientId)
	if err != nil {
		RedirectToError(w, r, data.RedirectUri, "invalid_client_id", "The client_id given in the request is not valid")
		return
	}

	// We put the default redirect_uri if it's missing from the query
	if data.RedirectUri == "" {
		data.RedirectUri = client.RedirectUri
	}

	// We check that the user is correct and redirect to an error if he's missing
	if !auth.CheckUser(data.Email, data.Password) {
		RedirectToError(w, r, data.RedirectUri, "access_denied", "The user did not consent.")
		return
	}

	switch data.ResponseType {
	case ResponseTypeCode:
		authorizeCodeFlow(w, r, data)
		return
	case ResponseTypeImplicit:
		implicitFlow(w, r, data)
		return
	default:
		RedirectToError(w, r, data.RedirectUri, "invalid_response_type", "The response_type given in the request is not valid")
		return
	}
}
