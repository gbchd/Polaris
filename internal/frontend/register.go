package frontend

import "net/http"

func ServeRegisterPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web-dev/register/index.html")
}
