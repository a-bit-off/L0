package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"service/internal/http-server/model"

	"service/internal/cache"
	"service/internal/storage/postgres"
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
		const op = "handlers.add.AddOrderPage"

		orderInfo := r.URL.Query().Get("orderInfo")
		var showMessage bool
		if message != "" {
			showMessage = true
		}

		data := AddPageData{OrderInfo: orderInfo, ShowMessage: showMessage, Message: message}

		lp := filepath.Join("public", "html", "add.html")
		tmpl, err := template.ParseFiles(lp)
		if err != nil {
			log.Printf("%s: %s\n", op, err)
			http.Error(w, fmt.Sprintf("Internal Server Error: %s", err), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("%s: %s\n", op, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Println("Template add.html executed successful!")
	}
}

// POST
func AddOrder(storage *postgres.Storage, cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.add.AddOrder"

		r.ParseForm()
		orderInfo := r.FormValue("orderInfo")

		if orderInfo == "" {
			AddOrderPage("Fields must be filled!")(w, r)
			return
		}
		var order model.Model
		if err := json.Unmarshal([]byte(orderInfo), &order); err != nil {
			AddOrderPage("Order not a model!")(w, r)
			return
		}

		// Add to storage
		err := storage.AddOrder(order)
		if err != nil {
			log.Printf("%s: %s\n", op, err)
			AddOrderPage(err.Error())(w, r)
			return
		}
		log.Println("Order added to db successfully!")

		// Add to cache
		byt, err := json.Marshal(order)
		if err != nil {
			log.Printf("%s: %w", op, err)
		} else {
			cache.SetDefault(order.OrderUID, byt)
		}

		log.Println("Order added to cache successfully!")

		AddOrderPage("Order added successfully!")(w, r)
	}
}
