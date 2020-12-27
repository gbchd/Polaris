package oauth

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/guillaumebchd/polaris/pkg/code"
)

func authorizeCodeFlow(w http.ResponseWriter, r *http.Request, data LoginFormData) {
	c_data := code.CodeData{
		Openid:     strings.Contains(data.Scope, "openid"),
		PKCEMethod: data.CodeChallengeMethod,
		PKCEValue:  data.CodeChallenge,
	}

	c, err := code.GenerateCode(c_data)
	if err != nil {
		RedirectToError(w, r, data.RedirectUri, "internal_error", "Something went wrong when creating your code.")
		return
	}

	ur, err := url.Parse(data.RedirectUri)
	if err != nil {
		panic(err)
	}

	q := ur.Query()
	q.Add("code", c)
	if data.State != "" {
		q.Add("state", data.State)
	}

	ur.RawQuery = q.Encode()
	http.Redirect(w, r, ur.String(), http.StatusSeeOther)
}

func tokenCodeFlow(w http.ResponseWriter, r *http.Request) {

}
