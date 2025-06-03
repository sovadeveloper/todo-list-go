package task

type Task struct {
	ID        int    `db:"id"`
	Title     string `db:"title"`
	Completed bool   `db:"completed"`
}

func (task *Task) MarkDone() {
	task.Completed = true
}
