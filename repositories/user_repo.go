package repositories

import (
	"errors"

	"github.com/alialaa77/TODO-App/config"
	"github.com/alialaa77/TODO-App/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo() *UserRepo {
	return &UserRepo{db: config.DB}
}

func (r *UserRepo) Create(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepo) GetByUsername(username string, u *models.User) error {
	res := r.db.Where("username = ?", username).First(u)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}

func (r *UserRepo) GetByID(id uint, u *models.User) error {
	res := r.db.First(u, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}
