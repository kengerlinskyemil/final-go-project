package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now, err := time.Parse(dateLayout, r.URL.Query().Get("now"))
	if err != nil {
		http.Error(w, "Некорректная дата в параметре now", http.StatusBadRequest)
		return
	}

	repeat := r.URL.Query().Get("repeat")
	if repeat == "" {
		http.Error(w, "Некорректная дата в параметре repeat", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(dateLayout, r.URL.Query().Get("date"))
	if err != nil {
		http.Error(w, "Некорректная дата в параметере date", http.StatusBadRequest)
		return
	}

	nextDate, err := nextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte(nextDate)); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func nextDate(now time.Time, date time.Time, repeat string) (string, error) {
	parts := strings.Split(strings.TrimSpace(repeat), " ")
	if len(parts) == 0 {
		return "", errors.New("empty repeat rule")
	}

	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid 'd' rule format: expected 'd <number>'")
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("invalid number of days in 'd' rule")
		}

		if interval < 1 || interval > 400 {
			return "", errors.New("day interval must be between 1 and 400")
		}

		date = date.AddDate(0, 0, interval)
		for !afterNow(date, now) {
			date = date.AddDate(0, 0, interval)
		}

		return date.Format(dateLayout), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("the 'y' rule does not take parameters")
		}

		date = date.AddDate(1, 0, 0)
		for !afterNow(date, now) {
			date = date.AddDate(1, 0, 0)
		}

		return date.Format(dateLayout), nil
	default:
		return "", errors.New("unknown repeat rule: " + rule)
	}
}

func afterNow(date, now time.Time) bool {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return dateOnly.After(nowOnly)
}
