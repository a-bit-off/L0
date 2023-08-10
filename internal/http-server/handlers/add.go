package handlers

import (
	"L0/internal/storage/postgres"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type AddPageData struct {
	OrderID     string
	OrderInfo   string
	ShowMessage bool
	Message     string
}

// Создание нового кэша с дефолтными настройками

// GET
func AddOrderPage(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("orderID")
		orderInfo := r.URL.Query().Get("orderInfo")
		var showMessage bool
		if message != "" {
			showMessage = true
		}

		data := AddPageData{OrderID: orderID, OrderInfo: orderInfo, ShowMessage: showMessage, Message: message}

		lp := filepath.Join("public", "html", "add.html")
		tmpl, err := template.ParseFiles(lp)
		if err != nil {
			http.Error(w, fmt.Sprintf("Internal Server Error: %s", err), http.StatusInternalServerError)
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
func AddOrder(storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		orderID := r.FormValue("orderID")
		orderInfo := r.FormValue("orderInfo")

		if orderID == "" || orderInfo == "" {
			AddOrderPage("All fields must be filled!")(w, r)
			return
		}

		err := storage.AddOrder(orderID, orderInfo)
		if err != nil {

			log.Println(err)
			AddOrderPage(err.Error())(w, r)
			return
		}

		AddOrderPage("Order added successfully!")(w, r)
	}
}
