package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"L0/internal/http-server/model"
)

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
