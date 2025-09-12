package models

import "time"

type Todo struct {
	ID          int        `json:"id" db:"id"`
	Title       string     `json:"title" db:"title`
	Completed   bool       `json:"completed" db:"completed"`
	Category    string     `json:"category" db:"category"`
	Priority    string     `json:"priority" db:"priority"`
	CompletedAt *time.Time `json:"completedAt" db:"completed_at"`
	DueDate     *time.Time `json:"dueDate" db:"due_date"`
}
