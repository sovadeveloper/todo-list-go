package main

import (
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"time"
	"todo-list/internal/api"
	"todo-list/internal/cache"
	"todo-list/internal/db"
	"todo-list/internal/migration"
	"todo-list/internal/task"
)

func main() {
	err := migration.Run()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func(database *sqlx.DB) {
		err := database.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(database)

	repo := task.NewRepository(database)

	loader := func() ([]task.Task, error) {
		return repo.List()
	}
	taskCache := cache.NewTaskCache(5*time.Second, loader)
	_ = taskCache.Init()
	handler := &api.Handler{Repo: repo, Cache: taskCache}

	log.Println("Server running on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler.Routes()))
}
