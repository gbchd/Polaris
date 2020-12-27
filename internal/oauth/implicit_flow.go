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

	at, err := token.CreateAccessToken(data.ClientId, data.Email)
	if err != nil {
		fmt.Println(err)
	}

	q.Add("access_token", at)
	q.Add("token_type", "Bearer")

	dur := int(token.LifetimeAccessToken.Seconds())
	q.Add("expires_in", strconv.Itoa(dur))

	if strings.Contains(data.Scope, "openid") {
		user, err := authentication.GetUser(data.Email)
		if err != nil {
			fmt.Println(err)
		}

		id_t, err := token.CreateIdToken(data.ClientId, user)
		q.Add("id_token", id_t)
	}

	ur.RawQuery = q.Encode()
	http.Redirect(w, r, ur.String(), http.StatusSeeOther)
}
