package oauth

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/guillaumebchd/polaris/pkg/authentication"
	"github.com/guillaumebchd/polaris/pkg/token"
)

func implicitFlow(w http.ResponseWriter, r *http.Request, data LoginFormData) {

	ur, err := url.Parse(data.RedirectUri)
	if err != nil {
		panic(err)
	}

	q := ur.Query()

	accessToken := token.AccessToken{
		Issuer:   "Polaris",
		Audience: data.ClientId,
	}
	at, err := accessToken.Encode()
	if err != nil {
		fmt.Println(err)
	}

	q.Add("access_token", at)
	q.Add("token_type", "Bearer")

	dur := int(token.AccessTokenLifetime.Seconds())
	q.Add("expires_in", strconv.Itoa(dur))

	if strings.Contains(data.Scope, "openid") {
		user, err := authentication.GetUser(data.Email)
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

		q.Add("id_token", id_t)
	}

	ur.RawQuery = q.Encode()
	http.Redirect(w, r, ur.String(), http.StatusSeeOther)
}
