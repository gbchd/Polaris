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
	at, err := token.CreateAccessToken(data.ClientId)
	if err != nil {
		fmt.Println(err)
	}

	res.TokenType = "Bearer"
	res.AccessToken = at
	res.ExpiresIn = strconv.Itoa(int(token.LifetimeAccessToken.Seconds()))

	if data.Scope != "" {
		res.Scope = data.Scope
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(resJson)
}
