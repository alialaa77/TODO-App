package models

import "time"

type Todo struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string     `json:"title" gorm:"size:255;not null"`
	Completed   bool       `json:"completed" gorm:"not null;default:false"`
	Category    string     `json:"category" gorm:"size:50;not null"`
	Priority    string     `json:"priority" gorm:"size:10;not null"`
	CompletedAt *time.Time `json:"completedAt" gorm:"column:completed_at"`
	DueDate     *time.Time `json:"dueDate" gorm:"column:due_date"`
}
