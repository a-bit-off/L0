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
	OrderID string
}

// GET
func FindOrderByIDPage(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("orderID")
	data := FindPageData{OrderID: orderID}

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

// POST
func FindOrderByID(storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		orderID := r.FormValue("orderID")

		if orderID == "" {
			FindOrderByIDPage(w, r) // Повторно отображаем страницу с предупреждением
			return
		}

		jsonB, err := storage.GetById(orderID)
		if err != nil {
			log.Println(fmt.Errorf("Error getting data from database: %s", err))
			http.Error(w, "Error getting data from database", http.StatusInternalServerError)
			return
		}

		if jsonB == nil {
			http.Error(w, "Record not found", http.StatusNotFound)
			return
		}

		OrderDetailsPage(jsonB)(w, r)
	}
}
