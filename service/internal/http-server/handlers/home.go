package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.home.HomePage"

	lp := filepath.Join("public", "html", "home.html")
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		log.Printf("%s: %s\n", op, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("%s: %s\n", op, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("Template home.html executed successful!")
}
