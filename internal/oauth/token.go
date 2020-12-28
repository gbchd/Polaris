package oauth

import (
	"encoding/base64"
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
	Code         string `schema:"code"`
	CodeVerifier string `schema:"code_verifier"`
	Scope        string `schema:"scope"`
	Username     string `schema:"username"` // Should be the email
	Password     string `schema:"password"`
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

	b64clientId, b64clientSecret, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Authorization Required"`)
		w.WriteHeader(401)
		return
	}

	clientId, clientSecret, err := decodeClientIdAndSecret(b64clientId, b64clientSecret)
	if err != nil {
		respondError(w, r, "error_decode", "Couldn't decode the base64 encoded client_id and client_secret")
		return
	}

	client, err := auth.GetClient(string(clientId))
	if err != nil || client.ClientSecret != clientSecret {
		respondError(w, r, "invalid_client_creds", "The client_id and client_secret given in the request do not match")
		return
	}

	data.ClientId = clientId

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

func decodeClientIdAndSecret(b64clientId string, b64ClientSecret string) (string, string, error) {
	clientId, err := base64.StdEncoding.DecodeString(b64clientId)
	if err != nil {
		return "", "", err
	}

	clientSecret, err := base64.StdEncoding.DecodeString(b64ClientSecret)
	if err != nil {
		return string(clientId), "", err
	}

	return string(clientId), string(clientSecret), err
}
