package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/guillaumebchd/polaris/pkg/authentication"
	auth "github.com/guillaumebchd/polaris/pkg/authentication"
	"github.com/guillaumebchd/polaris/pkg/token"
)

func resourceOwnerFlow(w http.ResponseWriter, r *http.Request, data TokenData) {

	// We check the username and password
	if !auth.CheckUser(data.Username, data.Password) {
		respondError(w, r, "access_denied", "The credentials are incorrect.")
		return
	}

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

	if strings.Contains(data.Scope, "openid") {
		user, err := authentication.GetUser(data.Username)
		if err != nil {
			fmt.Println(err)
		}

		idToken := token.IdToken{
			Issuer:   "Polaris",
			Audience: data.ClientId,
			Email:    user.Email,
			Name:     user.Name,
		}
		id_t, err := idToken.Encode()
		if err != nil {
			fmt.Println(err)
		}

		res.IdToken = id_t
		res.Scope = "openid"
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(resJson)
}
