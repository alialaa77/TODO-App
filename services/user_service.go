package services

import (
	"errors"
	"time"
	"your-project/models"
	"your-project/repositories"
)

func ValidateTodo(todo models.Todo) error {
	if todo.Title == "" {
		return errors.New("Title cannot be empty")
	}
	if todo.Priority != "Low" && todo.Priority != "Medium" && todo.Priority != "High" {
		return errors.New("Priority must be Low, Medium, or High")
	}
	if todo.DueDate != nil && todo.DueDate.Before(time.Now().UTC()) {
		return errors.New("Due date cannot be in the past")
	}
	return nil
}

func CreateTodoService(todo models.Todo) (models.Todo, error) {
	if err := ValidateTodo(todo); err != nil {
		return todo, err
	}
	if todo.Completed {
		now := time.Now().UTC()
		todo.CompletedAt = &now
	}
	id, err := repositories.CreateTodo(todo)
	if err != nil {
		return todo, err
	}
	todo.ID = id
	return todo, nil
}
