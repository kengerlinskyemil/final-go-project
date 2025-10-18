package api

import (
	"net/http"
	"time"

	"github.com/kengerlinskyemil/final-go-project/pkg/utils"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse("20060102", r.URL.Query().Get("now"))
	if err != nil {
		http.Error(w, "Некорректная дата в параметре now", http.StatusBadRequest)
		return
	}

	repeat := r.URL.Query().Get("repeat")
	if repeat == "" {
		http.Error(w, "Некорректная дата в параметре repeat", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("20060102", r.URL.Query().Get("date"))
	if err != nil {
		http.Error(w, "Некорректная дата в параметере date", http.StatusBadRequest)
		return
	}

	nextDate, err := utils.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(nextDate))
}
