// handlers/errors.go
package groupie

import (
	"html/template"
	"net/http"
)

type ErrorData struct {
	Code   int
	Errors []string
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, code int, errors []string) {
	w.WriteHeader(code)
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := ErrorData{Code: code, Errors: errors}
	tmpl.Execute(w, data)
}
