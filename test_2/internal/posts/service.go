package posts

import (
	"rootwritter/majoo_test_2_api/internal/models"
	"strconv"

	"gorm.io/gorm"
)

type Service interface {
	CreateNewPost(title, content string, userID uint) (*models.Post, error)
	GetAllPosts(page, limit string) ([]*models.Post, error)
	GetPostByID(id string) (*models.Post, error)
	UpdatePost(id string, title, content *string, userID uint) (*models.Post, error)
	DeletePost(id string, userID uint) error

	// Transaction support
	WithTransaction(tx *gorm.DB) Service
}

type service struct {
	repo Repository
	db   *gorm.DB // Store the original db instance for transactions
}

func NewService(repo Repository, db *gorm.DB) Service {
	return &service{repo: repo, db: db}
}

func (s *service) WithTransaction(tx *gorm.DB) Service {
	// Get a new repository instance with the transaction
	txRepo := s.repo.WithTransaction(tx)
	return &service{repo: txRepo, db: tx}
}

func (s *service) CreateNewPost(title, content string, userID uint) (*models.Post, error) {
	post := &models.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	// Anda bisa menambahkan logika bisnis di sini (misal: filter kata kasar)
	err := s.repo.CreatePost(post)
	return post, err
}

func (s *service) GetAllPosts(page, limit string) ([]*models.Post, error) {
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 10
	}

	offset := (pageInt - 1) * limitInt
	posts, err := s.repo.GetAllPosts(offset, limitInt)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *service) GetPostByID(id string) (*models.Post, error) {
	idInt, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}

	post, err := s.repo.GetPostByID(uint(idInt))
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *service) UpdatePost(id string, title, content *string, userID uint) (*models.Post, error) {
	idInt, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}

	// Check if the post belongs to the user
	existingPost, err := s.repo.GetPostByID(uint(idInt))
	if err != nil {
		return nil, err
	}

	if existingPost.UserID != userID {
		return nil, err // unauthorized
	}

	updateData := map[string]interface{}{}
	if title != nil {
		updateData["title"] = *title
	}
	if content != nil {
		updateData["content"] = *content
	}

	updatedPost, err := s.repo.UpdatePost(uint(idInt), updateData)
	if err != nil {
		return nil, err
	}

	return updatedPost, nil
}

func (s *service) DeletePost(id string, userID uint) error {
	idInt, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}

	// Check if the post belongs to the user
	existingPost, err := s.repo.GetPostByID(uint(idInt))
	if err != nil {
		return err
	}

	if existingPost.UserID != userID {
		return err // unauthorized
	}

	return s.repo.DeletePost(uint(idInt))
}
