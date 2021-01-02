package frontend

import (
	"fmt"
	"log"
	"net/http"

	"github.com/guillaumebchd/polaris/pkg/reset"
)

type RecoverFormData struct {
	Email string `schema:"email"`
}

func ServeRecoverPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/recover/index.html")
}

func RecoverFormHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var data RecoverFormData
	err = decoder.Decode(&data, r.Form)
	if err != nil {
		fmt.Println("Error in GET parameters : ", err)
	}

	code, err := reset.Generate(data.Email)
	if err != nil {
		fmt.Println("CODE GENERATION ERROR ", err)
	}

	// Change this to dns url
	link := "localhost:8000/reset/" + code
	fmt.Println(link)

	http.ServeFile(w, r, "web/recover/emailsent.html")
}
