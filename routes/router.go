package routes

import (
	"github.com/gin-gonic/gin"
	"todo.mod/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/todos", handlers.GetTodosHandler)
	r.POST("/todos", handlers.CreateTodoHandler)

	// Later: GET by ID, PUT, DELETE, filters, bulk update, etc.

	return r
}
