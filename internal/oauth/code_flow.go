package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/guillaumebchd/polaris/pkg/authentication"
	"github.com/guillaumebchd/polaris/pkg/code"
	"github.com/guillaumebchd/polaris/pkg/token"
)

func authorizeCodeFlow(w http.ResponseWriter, r *http.Request, data LoginFormData) {
	c_data := code.Data{
		Email:           data.Email,
		ChallengeMethod: data.CodeChallengeMethod,
		Challenge:       data.CodeChallenge,
	}

	c, err := code.Generate(c_data)
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

func tokenCodeFlow(w http.ResponseWriter, r *http.Request, data TokenData) {

	if data.Code == "" {
		respondError(w, r, "missing_code", "There is no code in the request")
		return
	}

	c, err := code.Get(data.Code)
	if err != nil {
		respondError(w, r, "invalid_code", "The given code is invalid or has expired")
		return
	}

	// PKCE - We check the code_challenge
	if data.CodeVerifier != "" {
		if c.ChallengeMethod == "S256" {
			hasher := sha256.New()
			hasher.Write([]byte(data.CodeVerifier))
			sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
			if sha != c.Challenge {
				respondError(w, r, "invalid_code_challenge", "The code_challenge and the code_verifier doesn't match.")
				return
			}
		} else if c.ChallengeMethod == "plaintext" {
			if data.CodeVerifier != c.Challenge {
				respondError(w, r, "invalid_code_challenge", "The code_challenge and the code_verifier doesn't match.")
				return
			}
		} else {
			respondError(w, r, "invalid_code_method", "The only acceptables methods are 'S256' and 'plaintext'.")
			return
		}
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
		user, err := authentication.GetUser(c.Email)
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

	refreshToken := token.RefreshToken{
		Issuer:   "Polaris",
		Audience: data.ClientId,
	}

	rt, err := refreshToken.Encode(at)
	if err != nil {
		fmt.Println(err)
	}

	res.RefreshToken = rt

	resJson, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(resJson)
}
