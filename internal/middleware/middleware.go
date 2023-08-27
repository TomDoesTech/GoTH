package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func RenderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateData := map[string]interface{}{
			"Title": "Default Title",
		}

		r = r.WithContext(context.WithValue(r.Context(), "templateData", templateData))

		next.ServeHTTP(w, r)
	})
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
