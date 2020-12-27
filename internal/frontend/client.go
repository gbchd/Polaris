package frontend

import "net/http"

func ServeClientPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web-dev/client/index.html")
}
