package repositories

import (
	"database/sql"

	"your-project/config"

	"your-project/models"
)

func GetAllTodos() ([]models.Todo, error) {
	rows, err := config.DB.Query("SELECT id, title, completed, category, priority, completed_at, due_date FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var t models.Todo
		var completedAt, dueDate sql.NullTime
		err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.Category, &t.Priority, &completedAt, &dueDate)
		if err != nil {
			return nil, err
		}
		if completedAt.Valid {
			t.CompletedAt = &completedAt.Time
		}
		if dueDate.Valid {
			t.DueDate = &dueDate.Time
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func CreateTodo(todo models.Todo) (int, error) {
	var id int
	err := config.DB.QueryRow(`
		INSERT INTO todos (title, completed, category, priority, completed_at, due_date)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		todo.Title, todo.Completed, todo.Category, todo.Priority, todo.CompletedAt, todo.DueDate,
	).Scan(&id)
	return id, err
}

// You will add: GetByID, Update, Delete, FilterByCategory, FilterByStatus, Search, BulkUpdate, etc.
