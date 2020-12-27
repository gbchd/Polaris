package frontend

import (
	"fmt"
	"log"
	"net/http"

	"github.com/guillaumebchd/polaris/internal/oauth"
	"github.com/guillaumebchd/polaris/pkg/authentication"
)

type RegisterFormData struct {
	Name      string `schema:"name"`
	Email     string `schema:"email"`
	Password  string `schema:"password"`
	VPassword string `schema:"vpassword"`
}

func ServeRegisterPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/register/index.html")
}

func RegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var data RegisterFormData
	err = decoder.Decode(&data, r.Form)
	if err != nil {
		fmt.Println("Error in GET parameters : ", err)
	}

	if data.Password != data.VPassword {
		oauth.RedirectToError(w, r, "", "invalid_password_match", "The two given passwords don't match.")
		return
	}

	_, err = authentication.CreateUser(data.Name, data.Email, data.Password)
	if err != nil {
		oauth.RedirectToError(w, r, "", "internal_error", "Something went wrong when creating the user.")
		return
	}

	http.ServeFile(w, r, "web/register/valid.html")

}
