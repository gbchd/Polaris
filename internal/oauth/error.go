package oauth

import (
	"net/http"
	"net/url"
)

func RedirectToError(w http.ResponseWriter, r *http.Request, redirect_url string, err_code string, err_desc string) {
	if redirect_url == "" {
		redirect_url = "/error"
	}

	ur, err := url.Parse(redirect_url)
	if err != nil {
		panic(err)
	}

	q := ur.Query()
	q.Add("error", err_code)
	q.Add("error_description", err_desc)
	ur.RawQuery = q.Encode()

	http.Redirect(w, r, ur.String(), http.StatusSeeOther)
}
