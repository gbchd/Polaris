package oauth

import (
	"fmt"
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

var decoder = schema.NewDecoder()

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	var data AuthorizeData
	err := decoder.Decode(&data, r.URL.Query())
	if err != nil {
		fmt.Println("Error in GET parameters : ", err)
	}

	// We check that the client id is here and is not invalid
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

	switch data.ResponseType {
	case ResponseTypeCode:
		if data.CodeChallengeMethod == "" || data.CodeChallenge == "" {
			authorizationCodeFlow(w, r, data)
		} else {
			authorizationCodeFlowPKCE(w, r, data)
		}
		return
	case ResponseTypeImplicit:
		implicitFlow(w, r, data)
		return
	default:
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
}

/*
	FLOWS
*/

func authorizationCodeFlow(w http.ResponseWriter, r *http.Request, data AuthorizeData) {

}

func authorizationCodeFlowPKCE(w http.ResponseWriter, r *http.Request, data AuthorizeData) {

}

func implicitFlow(w http.ResponseWriter, r *http.Request, data AuthorizeData) {

}
