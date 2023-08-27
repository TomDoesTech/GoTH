package web

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/tomdoestech/goth/internal/middleware"
)

type WebHTTPParams struct {
	WebHandler *WebHandler
	Mux        *chi.Mux
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func RenderTemplate(w http.ResponseWriter, tmplName string, data interface{}, r *http.Request) {

	_, user, _ := jwtauth.FromContext(r.Context())

	tmpl, err := template.ParseFiles(
		"templates/"+tmplName,
		"templates/partial/header.html",
		"templates/partial/nav.html",
		"templates/partial/footer.html",
		"templates/partial/base.html",
	)

	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	data.(map[string]interface{})["User"] = user

	scriptNonce := "htmx_" + generateRandomString(8)
	styleNonce := "tw_" + generateRandomString(8)

	data.(map[string]interface{})["scriptNonce"] = scriptNonce
	data.(map[string]interface{})["styleNonce"] = styleNonce

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func NewWebHTTP(p WebHTTPParams) {
	r := p.Mux

	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title": "Login",
		}

		RenderTemplate(w, "login.html", data, r)
	})

	r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title": "Register",
		}

		RenderTemplate(w, "register.html", data, r)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		data := map[string]interface{}{
			"Title": "My Website",
		}

		// Render the home.html template and inject data
		middleware.RenderTemplate(w, "home.html", data, r)
	})

	r.Get("/about", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title": "About",
		}

		// Render the home.html template and inject data
		middleware.RenderTemplate(w, "about.html", data, r)
	})

}
