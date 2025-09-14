package routes

import (
	"github.com/alialaa77/TODO-App/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	h := handlers.NewTodoHandler()
	h.RegisterRoutes(r)
	return r
}
