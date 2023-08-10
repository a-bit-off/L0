package handlers

import (
	"L0/internal/storage/postgres"
	"html/template"
	"net/http"
	"path/filepath"
)

type AddPageData struct {
	OrderID   string
	OrderInfo string
}

// GET
func AddOrderPage(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("orderID")
	orderInfo := r.URL.Query().Get("orderInfo")
	data := AddPageData{OrderID: orderID, OrderInfo: orderInfo}

	lp := filepath.Join("public", "html", "add.html")
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
func AddOrder(storage *postgres.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		orderID := r.FormValue("orderID")
		orderInfo := r.FormValue("orderInfo")

		if orderID == "" || orderInfo == "" {
			AddOrderPage(w, r) // TODO: Повторно отображаем страницу с предупреждением
			return
		}

		err := storage.AddOrder(orderID, orderInfo)
		if err != nil {
			http.Error(w, "Error with add order", http.StatusInternalServerError)
			return
		}

		// TODO: Дать понять что все ок!
	}
}
