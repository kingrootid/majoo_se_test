package comments

import (
	"rootwritter/majoo_test_2_api/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	CreateComment(comment *models.Comment) error
	GetCommentByID(id uint) (*models.Comment, error)
	GetCommentsByPostID(postID uint, offset, limit int) ([]*models.Comment, error)
	GetCommentsByUserID(userID uint, offset, limit int) ([]*models.Comment, error)
	UpdateComment(id uint, data map[string]interface{}) (*models.Comment, error)
	DeleteComment(id uint) error

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

func (r *repository) CreateComment(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *repository) GetCommentByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.First(&comment, id).Error
	return &comment, err
}

func (r *repository) GetCommentsByPostID(postID uint, offset, limit int) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := r.db.Offset(offset).Limit(limit).Where("post_id = ?", postID).Preload("User").Preload("Post").Find(&comments).Error
	return comments, err
}

func (r *repository) GetCommentsByUserID(userID uint, offset, limit int) ([]*models.Comment, error) {
	var comments []*models.Comment
	err := r.db.Offset(offset).Limit(limit).Where("user_id = ?", userID).Preload("Post").Preload("User").Find(&comments).Error
	return comments, err
}

func (r *repository) UpdateComment(id uint, data map[string]interface{}) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&comment).Updates(data).Error
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *repository) DeleteComment(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}