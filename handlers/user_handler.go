package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alialaa77/TODO-App/models"
	"github.com/alialaa77/TODO-App/repositories"
	"github.com/alialaa77/TODO-App/services"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	repo    *repositories.TodoRepo
	service *services.TodoService
}

func NewTodoHandler() *TodoHandler {
	r := repositories.NewTodoRepo()
	return &TodoHandler{
		repo:    r,
		service: services.NewTodoService(r),
	}
}

type todoInput struct {
	Title     string  `json:"title"`
	Completed bool    `json:"completed"`
	Category  string  `json:"category"`
	Priority  string  `json:"priority"`
	DueDate   *string `json:"dueDate"`
}

func (h *TodoHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/todos", h.GetAll)
	r.GET("/todos/:id", h.GetByID)
	r.GET("/todos/category/:category", h.GetByCategory)
	r.GET("/todos/status/:status", h.GetByStatus)
	r.GET("/todos/search", h.Search)
	r.POST("/todos", h.Create)
	r.PUT("/todos/:id", h.Update)
	r.PUT("/todos/category/:category", h.BulkUpdateCategory)
	r.DELETE("/todos/:id", h.DeleteByID)
	r.DELETE("/todos", h.DeleteAll)
}

func (h *TodoHandler) GetAll(c *gin.Context) {
	var todos []models.Todo
	if err := h.repo.GetAll(&todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var todo models.Todo
	if err := h.repo.GetByID(uint(id), &todo); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) GetByCategory(c *gin.Context) {
	category := c.Param("category")
	var todos []models.Todo
	if err := h.repo.GetByCategory(category, &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) GetByStatus(c *gin.Context) {
	s := strings.ToLower(c.Param("status"))
	completed := s == "true" || s == "1"
	var todos []models.Todo
	if err := h.repo.GetByStatus(completed, &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) Search(c *gin.Context) {
	q := strings.ToLower(c.Query("q"))
	var todos []models.Todo
	if err := h.repo.SearchByTitle(q, &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) Create(c *gin.Context) {
	var in todoInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var due *time.Time
	if in.DueDate != nil && *in.DueDate != "" {
		t, err := time.Parse(time.RFC3339, *in.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dueDate format"})
			return
		}
		t = t.UTC()

		ref := time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC)
		if t.Before(ref) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Due date cannot be in the past"})
			return
		}
		due = &t
	}
	todo := &models.Todo{
		Title:     in.Title,
		Completed: in.Completed,
		Category:  in.Category,
		Priority:  in.Priority,
		DueDate:   due,
	}
	if err := h.service.CreateTodo(todo); err != nil {
		code := http.StatusBadRequest
		if err.Error() == "priority must be Low, Medium, or High" {
			code = http.StatusBadRequest
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var existing models.Todo
	if err := h.repo.GetByID(uint(id), &existing); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	var in todoInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	var due *time.Time
	if in.DueDate != nil && *in.DueDate != "" {
		t, err := time.Parse(time.RFC3339, *in.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dueDate format"})
			return
		}
		t = t.UTC()
		ref := time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC)
		if t.Before(ref) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Due date cannot be in the past"})
			return
		}
		due = &t
	}
	update := &models.Todo{
		Title:     in.Title,
		Completed: in.Completed,
		Category:  in.Category,
		Priority:  in.Priority,
		DueDate:   due,
	}
	if err := h.service.UpdateTodo(&existing, update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, existing)
}

func (h *TodoHandler) BulkUpdateCategory(c *gin.Context) {
	category := c.Param("category")
	var payload struct {
		Completed bool `json:"completed"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	todos, err := h.repo.UpdateCategoryStatus(category, payload.Completed, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// set completedAt properly for those set to true
	if payload.Completed {
		now := time.Now().UTC()
		for i := range todos {
			todos[i].Completed = true
			todos[i].CompletedAt = &now
			_ = h.repo.Update(&todos[i])
		}
	}
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) DeleteByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.repo.DeleteByID(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}

func (h *TodoHandler) DeleteAll(c *gin.Context) {
	if err := h.repo.DeleteAll(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All todos deleted"})
}
