package task

import "github.com/jmoiron/sqlx"

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db}
}

func (repo *Repository) Add(title string) error {
	_, err := repo.db.Exec("INSERT INTO tasks (title) VALUES ($1)", title)
	return err
}

func (repo *Repository) List() ([]Task, error) {
	var tasks []Task
	err := repo.db.Select(&tasks, "SELECT * FROM tasks")
	return tasks, err
}

func (repo *Repository) Complete(id int) error {
	_, err := repo.db.Exec("UPDATE tasks SET completed = TRUE WHERE id = $1", id)
	return err
}

func (repo *Repository) Delete(id int) error {
	_, err := repo.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	return err
}
