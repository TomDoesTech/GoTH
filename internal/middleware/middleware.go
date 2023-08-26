package middleware

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
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

func RenderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl, err := template.ParseFiles(
		"templates/"+tmplName,
		"templates/partial/header.html",
		"templates/partial/footer.html",
		"templates/partial/base.html",
	)

	if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
