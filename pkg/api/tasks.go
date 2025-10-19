package api

import (
	"net/http"
	"strconv"

	"github.com/kengerlinskyemil/final-go-project/pkg/db"
)

type TasksResp struct {
	Tasks []map[string]string `json:"tasks"`
}

const tasksLimit = 50

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSONStatus(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}
	tasks, err := db.Tasks(tasksLimit)
	if err != nil {
		writeJSONStatus(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	out := make([]map[string]string, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, map[string]string{
			"id":      strconv.Itoa(int(t.ID)),
			"date":    t.Date,
			"title":   t.Title,
			"comment": t.Comment,
			"repeat":  t.Repeat,
		})
	}
	writeJSON(w, TasksResp{Tasks: out})
}
