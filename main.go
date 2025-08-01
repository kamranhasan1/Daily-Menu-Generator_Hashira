package main

import (
	"log"
	"net/http"
	"menu-api/api"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next(w, r)
	}
}

func main() {
	http.HandleFunc("/menu", enableCORS(api.MenuHandler))
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}