package oauth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	auth "github.com/guillaumebchd/polaris/pkg/authentication"
)

const (
	GrantTypeCode                  = "authorization_code"
	GrantTypeClientCred            = "client_credentials"
	GrantTypeResourceOwnerPassword = "password"
)

type TokenData struct {
	GrantType    string `schema:"grant_type"`
	ClientId     string `schema:"client_id"`
	ClientSecret string `schema:"client_secret"`
	Code         string `schema:"code"`
	CodeVerifier string `schema:"code_verifier"`
	Scope        string `schema:"scope"`
}

type ErrorResponseJSON struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type ResponseJSON struct {
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    string `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	Scope        string `json:"scope,omitempty"`
	IdToken      string `json:"id_token,omitempty"`
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var data TokenData
	err = decoder.Decode(&data, r.Form)
	if err != nil {
		fmt.Println("Error in GET parameters : ", err)
	}

	// We check that the client id is here and is not invalid
	if data.ClientId == "" {
		respondError(w, r, "missing_client_id", "The client_id is missing from the request")
		return
	}
	_, err = auth.GetClient(data.ClientId)
	if err != nil {
		respondError(w, r, "invalid_client_id", "The client_id given in the request is not valid")
		return
	}

	switch data.GrantType {
	case GrantTypeCode:
		tokenCodeFlow(w, r, data)
		return
	case GrantTypeClientCred:
		clientsCredsFlow(w, r, data)
		return
	case GrantTypeResourceOwnerPassword:
		resourceOwnerFlow(w, r, data)
		return
	default:
		respondError(w, r, "invalid_grant_type", "The grant_type given in the request is not valid")
		return
	}

}

func respondError(w http.ResponseWriter, r *http.Request, err_code string, err_desc string) {
	e := ErrorResponseJSON{
		Error:            err_code,
		ErrorDescription: err_desc,
	}
	j, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write(j)
	return
}
