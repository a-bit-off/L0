package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func HomePage(w http.ResponseWriter, r *http.Request) {

	lp := filepath.Join("public", "html", "home.html")
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
