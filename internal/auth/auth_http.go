package auth

import (
	"github.com/go-chi/chi/v5"
)

type AuthHTTPParams struct {
	AuthHandler *AuthHandler
	Mux         *chi.Mux
}

func NewAuthHTTP(p AuthHTTPParams) {

	r := p.Mux

	r.Post("/api/login", p.AuthHandler.Login)

	r.Post("/api/register", p.AuthHandler.Register)

	r.Post("/api/logout", p.AuthHandler.Logout)
}
