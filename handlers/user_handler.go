package handlers

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"

	"  github.com/alialaa77/TODO-App/models"
	"github.com/alialaa77/TODO-App/repositories"
	"github.com/alialaa77/TODO-App/utils"
)

type AuthHandler struct {
	repo *repositories.UserRepo
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{repo: repositories.NewUserRepo()}
}

type authInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(r *gin.Engine) {
	r.POST("/signup", h.Signup)
	r.POST("/login", h.Login)
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var in authInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var existing models.User
	if err := h.repo.GetByUsername(in.Username, &existing); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already taken"})
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	user := models.User{
		Username: in.Username,
		Password: string(hashed),
		Role:     "user",
	}
	if err := h.repo.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create user"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var in authInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	var user models.User
	if err := h.repo.GetByUsername(in.Username, &user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
