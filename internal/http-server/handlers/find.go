package handlers

import (
	"L0/internal/storage/postgres"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type FindPageData struct {
	OrderID     string
	ShowMessage bool
	Message     string
}

// GET
func FindOrderByIDPage(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("orderID")
		var showMessage bool
		if message != "" {
			showMessage = true
		}
		data := FindPageData{OrderID: orderID, ShowMessage: showMessage, Message: message}

		lp := filepath.Join("public", "html", "find.html")
		tmpl, err := template.ParseFiles(lp)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

// POST
func FindOrderByID(storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		orderID := r.FormValue("orderID")

		if orderID == "" {
			FindOrderByIDPage("Field must be filled!")(w, r) // Повторно отображаем страницу с предупреждением
			return
		}

		jsonB, err := storage.GetById(orderID)
		if err != nil {
			log.Println(fmt.Errorf("Error getting data from database: %s", err))
			http.Error(w, "Error getting data from database", http.StatusInternalServerError)
			return
		}

		if jsonB == nil {
			FindOrderByIDPage("Nothing found for this id!")(w, r) // Повторно отображаем страницу с предупреждением
			return
		}

		OrderDetailsPage(jsonB)(w, r)
	}
}
