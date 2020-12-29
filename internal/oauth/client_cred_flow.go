package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/guillaumebchd/polaris/pkg/token"
)

func clientsCredsFlow(w http.ResponseWriter, r *http.Request, data TokenData) {
	// We should handle the scopes here
	res := ResponseJSON{}

	accessToken := token.AccessToken{
		Issuer:   "Polaris",
		Audience: data.ClientId,
	}
	at, err := accessToken.Encode()
	if err != nil {
		fmt.Println(err)
	}

	res.TokenType = "Bearer"
	res.AccessToken = at
	res.ExpiresIn = strconv.Itoa(int(token.AccessTokenLifetime.Seconds()))

	if data.Scope != "" {
		res.Scope = data.Scope
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(resJson)
}
