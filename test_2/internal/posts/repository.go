package posts

import (
	"rootwritter/majoo_test_2_api/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	CreatePost(post *models.Post) error
	GetPostByID(id uint) (*models.Post, error)
	GetAllPosts(offset, limit int) ([]*models.Post, error)
	UpdatePost(id uint, data map[string]interface{}) (*models.Post, error)
	DeletePost(id uint) error

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

func (r *repository) CreatePost(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *repository) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("Comments").Preload("Comments.User").Preload("Comments.Post").First(&post, id).Error
	return &post, err
}

func (r *repository) GetAllPosts(offset, limit int) ([]*models.Post, error) {
	var posts []*models.Post
	err := r.db.Offset(offset).Limit(limit).Preload("Comments").Preload("Comments.User").Preload("Comments.Post").Find(&posts).Error
	return posts, err
}

func (r *repository) UpdatePost(id uint, data map[string]interface{}) (*models.Post, error) {
	var post models.Post
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&post).Updates(data).Error
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *repository) DeletePost(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}
