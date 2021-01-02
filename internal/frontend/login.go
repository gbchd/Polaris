package frontend

import (
	"net/http"
)

func ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/login/index.html")
}
