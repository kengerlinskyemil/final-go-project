package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/kengerlinskyemil/final-go-project/pkg/db"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Error writing JSON: %v", err)
	}
}

func writeJSONStatus(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Error writing JSON: %v", err)
	}
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		w.Header().Set("Allow", http.MethodPost+", "+http.MethodGet+", "+http.MethodPut+", "+http.MethodDelete)
		writeJSONStatus(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONStatus(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	if len(task.Title) == 0 {
		writeJSONStatus(w, http.StatusBadRequest, map[string]any{"error": "не указан заголовок задачи"})
		return
	}

	now := time.Now()
	if len(task.Date) == 0 {
		task.Date = now.Format(dateLayout)
	}

	t, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		writeJSONStatus(w, http.StatusBadRequest, map[string]any{"error": "некорректная дата"})
		return
	}

	var next string
	if len(task.Repeat) > 0 {
		next, err = nextDate(now, t, task.Repeat)
		if err != nil {
			writeJSONStatus(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
	}

	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(dateLayout)
		} else {
			task.Date = next
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONStatus(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]any{"id": id})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		writeJSONStatus(w, http.StatusBadRequest, map[string]any{"error": "не указан идентификатор"})
		return
	}
	t, err := db.GetTask(id)
	if err != nil {
		writeJSONStatus(w, http.StatusNotFound, map[string]any{"error": "задача не найдена"})
		return
	}
	writeJSON(w, map[string]string{
		"id":      strconv.FormatInt(t.ID, 10),
		"date":    t.Date,
		"title":   t.Title,
		"comment": t.Comment,
		"repeat":  t.Repeat,
	})
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var in map[string]any
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONStatus(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	sid, _ := in["id"].(string)
	if len(sid) == 0 {
		writeJSONStatus(w, http.StatusBadRequest, map[string]any{"error": "не указан идентификатор"})
		return
	}

	t := db.Task{}
	t.ID, _ = strconv.ParseInt(sid, 10, 64)
	if v, ok := in["date"].(string); ok {
		t.Date = v
	}
	if v, ok := in["title"].(string); ok {
		t.Title = v
	}
	if v, ok := in["comment"].(string); ok {
		t.Comment = v
	}
	if v, ok := in["repeat"].(string); ok {
		t.Repeat = v
	}

	if len(t.Title) == 0 {
		writeJSONStatus(w, http.StatusBadRequest, map[string]any{"error": "не указан заголовок задачи"})
		return
	}

	now := time.Now()
	if len(t.Date) == 0 {
		t.Date = now.Format(dateLayout)
	}
	dt, err := time.Parse(dateLayout, t.Date)
	if err != nil {
		writeJSONStatus(w, http.StatusBadRequest, map[string]any{"error": "некорректная дата"})
		return
	}
	var next string
	if len(t.Repeat) > 0 {
		next, err = nextDate(now, dt, t.Repeat)
		if err != nil {
			writeJSONStatus(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
	}
	if afterNow(now, dt) {
		if len(t.Repeat) == 0 {
			t.Date = now.Format(dateLayout)
		} else {
			t.Date = next
		}
	}

	if err := db.UpdateTask(&t); err != nil {
		writeJSONStatus(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]any{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		writeJSONStatus(w, http.StatusBadRequest, map[string]any{"error": "не указан идентификатор"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJSONStatus(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]any{})
}

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSONStatus(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		writeJSONStatus(w, http.StatusBadRequest, map[string]any{"error": "не указан идентификатор"})
		return
	}
	t, err := db.GetTask(id)
	if err != nil {
		writeJSONStatus(w, http.StatusNotFound, map[string]any{"error": "задача не найдена"})
		return
	}

	if len(t.Repeat) == 0 {
		if err := db.DeleteTask(id); err != nil {
			writeJSONStatus(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]any{})
		return
	}

	now := time.Now()
	cur, err := time.Parse(dateLayout, t.Date)
	if err != nil {
		writeJSONStatus(w, http.StatusBadRequest, map[string]any{"error": "некорректная дата"})
		return
	}
	next, err := nextDate(now, cur, t.Repeat)
	if err != nil {
		writeJSONStatus(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	if err := db.UpdateDate(next, id); err != nil {
		writeJSONStatus(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]any{})
}
