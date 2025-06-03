package main

import (
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"todo-list/internal/api"
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
	handler := &api.Handler{Repo: repo}

	log.Println("Server running on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler.Routes()))
}
