package routes

import (
	"github.com/alialaa77/TODO-App/handlers"
	"todo.mod/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := handlers.NewAuthHandler()
	auth.Register(r)

	h := handlers.NewTodoHandler()
	todos := r.Group("/todos")
	todos.Use(middlewares.JWTMiddleware())
	h.RegisterRoutes(todos)
	return r
}
