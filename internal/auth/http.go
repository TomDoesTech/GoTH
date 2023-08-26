package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tomdoestech/goth/internal/middleware"
)

type AuthHTTPParams struct {
	AuthHandler *AuthHandler
	Mux         *chi.Mux
}

func NewAuthHTTP(p AuthHTTPParams) {

	r := p.Mux

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title":  "Login",
			"Header": "Welcome to my website!",
			"Footer": "© 2023 My Website",
		}

		// Render the home.html template and inject data
		middleware.RenderTemplate(w, "login.html", data)
	})

	r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title": "Register",
		}

		// Render the home.html template and inject data
		middleware.RenderTemplate(w, "register.html", data)
	})

	r.Post("/api/login", p.AuthHandler.Login)

	r.Post("/api/register", p.AuthHandler.Register)
}
