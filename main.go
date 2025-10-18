package main

import (
	"log"
	"net/http"

	"github.com/kengerlinskyemil/final-go-project/pkg/api"
	"github.com/kengerlinskyemil/final-go-project/pkg/db"
)

var config = struct {
	Port   string
	WebDir string
}{
	Port:   "7540",
	WebDir: "web",
}

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.GetDB().Close()

	http.HandleFunc("/api/nextdate", api.NextDateHandler)
	http.HandleFunc("/api/task", api.TaskHandler)
	http.HandleFunc("/api/task/done", api.TaskDoneHandler)
	http.HandleFunc("/api/tasks", api.TasksHandler)

	http.Handle("/", http.FileServer(http.Dir(config.WebDir)))

	log.Printf("Сервер запущен на порту %s", config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
