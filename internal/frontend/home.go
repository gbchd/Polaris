package frontend

import "net/http"

func ServeHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/home/index.html")
}
