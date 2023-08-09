package handlers

import (
	"L0/internal/storage/postgres"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"L0/internal/http-server/model"
)

type PageData struct {
	OrderID string
}

// GET
func StartGetIDPage(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("orderID")
	data := PageData{OrderID: orderID}

	lp := filepath.Join("public", "html", "startPage.html")
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
func StartGetID(storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		orderID := r.FormValue("orderID")

		if orderID == "" {
			StartGetIDPage(w, r) // Повторно отображаем страницу с предупреждением
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

func OrderDetailsPage(jsonB []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data model.Model

		json.NewDecoder(strings.NewReader(string(jsonB))).Decode(&data)

		lp := filepath.Join("public", "html", "orderDetails.html")
		tmpl, err := template.ParseFiles(lp)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			fmt.Println("Template Execution Error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
