package repositories

import (
	"errors"

	"github.com/alialaa77/TODO-App/config"
	"github.com/alialaa77/TODO-App/models"

	"gorm.io/gorm"
)

type TodoRepo struct {
	db *gorm.DB
}

func NewTodoRepo() *TodoRepo {
	return &TodoRepo{db: config.DB}
}

func (r *TodoRepo) AutoMigrate() error {
	return r.db.AutoMigrate(&models.Todo{})
}

func (r *TodoRepo) GetAll(todos *[]models.Todo) error {
	return r.db.Find(todos).Error
}

func (r *TodoRepo) GetByID(id uint, todo *models.Todo) error {
	res := r.db.First(todo, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}

func (r *TodoRepo) GetByCategory(category string, todos *[]models.Todo) error {
	return r.db.Where("category = ?", category).Find(todos).Error
}

func (r *TodoRepo) GetByStatus(completed bool, todos *[]models.Todo) error {
	return r.db.Where("completed = ?", completed).Find(todos).Error
}

func (r *TodoRepo) SearchByTitle(q string, todos *[]models.Todo) error {
	return r.db.Where("LOWER(title) LIKE ?", "%"+q+"%").Find(todos).Error
}

func (r *TodoRepo) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

func (r *TodoRepo) Update(todo *models.Todo) error {
	return r.db.Save(todo).Error
}

func (r *TodoRepo) UpdateCategoryStatus(category string, completed bool, completedAt *string) ([]models.Todo, error) {
	var todos []models.Todo
	if err := r.db.Where("category = ?", category).Find(&todos).Error; err != nil {
		return nil, err
	}
	for i := range todos {
		todos[i].Completed = completed
		if completed {
			// set completedAt in caller (service) - repository keeps Save simple
		} else {
			todos[i].CompletedAt = nil
		}
		if err := r.db.Save(&todos[i]).Error; err != nil {
			return nil, err
		}
	}
	return todos, nil
}

func (r *TodoRepo) DeleteByID(id uint) error {
	res := r.db.Delete(&models.Todo{}, id)
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}

func (r *TodoRepo) DeleteAll() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Todo{}).Error
}
