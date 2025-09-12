package handlers

import (
	"net/http"
	"your-project/models"
	"your-project/services"

	"github.com/gin-gonic/gin"
)

func CreateTodoHandler(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if todo.DueDate != nil {
		utc := todo.DueDate.UTC()
		todo.DueDate = &utc
	}

	createdTodo, err := services.CreateTodoService(todo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdTodo)
}

func GetTodosHandler(c *gin.Context) {
	todos, err := services.GetAllTodosService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	c.JSON(http.StatusOK, todos)
}
