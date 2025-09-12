package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"githup.com/gin/gonic/gin"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos []Todo
var nextID = 1

func getTodos(context *gin.Context) {
	context.JSON(http.StatusOK, todos)
}

func addTodo(context *gin.Context) {

	var newTodo Todo

	if err := context.BindJSON(&newTodo); err != nil {
		return
	}

	newTodo.ID = nextID
	nextID++

	todos = append(todos, newTodo)

	context.IndentedJSON(http.StatusCreated, newTodo)
}

func main() {
	router := gin.Default() // gin,default is the function that create new gin router.

	// GET /todos -> return current slice
	router.GET("/todos", getTodos)

	router.POST("/todos", addTodo)

	// start server
	router.Run(":8080")
}
