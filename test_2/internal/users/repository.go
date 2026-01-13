package users

import (
	"rootwritter/majoo_test_2_api/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(id uint, data map[string]interface{}) (*models.User, error)
	DeleteUser(id uint) error

	// Transaction methods
	WithTransaction(tx *gorm.DB) Repository
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) WithTransaction(tx *gorm.DB) Repository {
	return &repository{db: tx}
}

func (r *repository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *repository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *repository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *repository) UpdateUser(id uint, data map[string]interface{}) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&user).Updates(data).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) DeleteUser(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}