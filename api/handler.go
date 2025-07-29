package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"menu-api/internal/logic"
)

func MenuHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, "Missing date parameter", http.StatusBadRequest)
		return
	}

	days := 3 
	if daysParam := r.URL.Query().Get("days"); daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && (d == 3 || d == 7) {
			days = d
		}
	}

	var menu []logic.DailyMenu
	var err error

	if days == 7 {
		menu, err = logic.GenerateWeeklyMenu(date)
	} else {
		menu, err = logic.GenerateThreeDayMenu(date)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(menu); err != nil {
		log.Println("Error encoding response:", err)
	}
}