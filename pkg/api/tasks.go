package api

import (
	"net/http"
	"strconv"

	"github.com/kengerlinskyemil/final-go-project/pkg/db"
)

type TasksResp struct {
	Tasks []map[string]string `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, map[string]any{"error": "unsupported method"})
		return
	}
	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
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
