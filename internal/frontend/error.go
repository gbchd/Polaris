package frontend

import (
	"html/template"
	"log"
	"net/http"
)

type ErrorData struct {
	Code        string `schema:"error" json:"error"`
	Description string `schema:"error_description" json:"error_description"`
}

func ErrorPageHandler(w http.ResponseWriter, r *http.Request) {
	var data ErrorData
	err := decoder.Decode(&data, r.URL.Query())
	if err != nil {
		log.Fatal("Error in GET parameters : ", err)
	}

	tmpl := template.Must(template.ParseFiles("web/error/index.html"))
	tmpl.Execute(w, data)
}
