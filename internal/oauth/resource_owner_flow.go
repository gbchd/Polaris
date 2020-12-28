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

	at, err := token.CreateAccessToken(data.ClientId)
	if err != nil {
		fmt.Println(err)
	}

	res.TokenType = "Bearer"
	res.AccessToken = at
	res.ExpiresIn = strconv.Itoa(int(token.LifetimeAccessToken.Seconds()))

	if strings.Contains(data.Scope, "openid") {
		user, err := authentication.GetUser(data.Username)
		if err != nil {
			fmt.Println(err)
		}

		id_t, err := token.CreateIdToken(data.ClientId, user)
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
