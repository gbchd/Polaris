package token

import "net/http"

func ServePubKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(pubKey.Raw)
}
