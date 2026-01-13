package comments

import (
	"errors"
	"strconv"

	"rootwritter/majoo_test_2_api/internal/models"
)

type Service interface {
	CreateNewComment(content string, postID, userID uint) (*models.Comment, error)
	GetCommentByID(id string) (*models.Comment, error)
	GetCommentsByPostID(postID string, page, limit string) ([]*models.Comment, error)
	GetCommentsByUserID(userID string, page, limit string) ([]*models.Comment, error)
	UpdateComment(id string, content *string, userID uint) (*models.Comment, error)
	DeleteComment(id string, userID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateNewComment(content string, postID, userID uint) (*models.Comment, error) {
	comment := &models.Comment{
		Content: content,
		PostID:  postID,
		UserID:  userID,
	}

	err := s.repo.CreateComment(comment)
	return comment, err
}

func (s *service) GetCommentByID(id string) (*models.Comment, error) {
	idInt, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}

	comment, err := s.repo.GetCommentByID(uint(idInt))
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *service) GetCommentsByPostID(postID string, page, limit string) ([]*models.Comment, error) {
	postIDInt, err := strconv.ParseUint(postID, 10, 32)
	if err != nil {
		return nil, err
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 10
	}

	offset := (pageInt - 1) * limitInt
	comments, err := s.repo.GetCommentsByPostID(uint(postIDInt), offset, limitInt)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *service) GetCommentsByUserID(userID string, page, limit string) ([]*models.Comment, error) {
	userIDInt, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, err
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 10
	}

	offset := (pageInt - 1) * limitInt
	comments, err := s.repo.GetCommentsByUserID(uint(userIDInt), offset, limitInt)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *service) UpdateComment(id string, content *string, userID uint) (*models.Comment, error) {
	idInt, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}

	// Check if the comment belongs to the user
	existingComment, err := s.repo.GetCommentByID(uint(idInt))
	if err != nil {
		return nil, err
	}

	if existingComment.UserID != userID {
		return nil, errors.New("unauthorized") // unauthorized
	}

	if content == nil {
		return nil, errors.New("content cannot be empty")
	}

	updateData := map[string]interface{}{
		"Content": *content,
	}

	updatedComment, err := s.repo.UpdateComment(uint(idInt), updateData)
	if err != nil {
		return nil, err
	}

	return updatedComment, nil
}

func (s *service) DeleteComment(id string, userID uint) error {
	idInt, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}

	// Check if the comment belongs to the user
	existingComment, err := s.repo.GetCommentByID(uint(idInt))
	if err != nil {
		return err
	}

	if existingComment.UserID != userID {
		return errors.New("unauthorized") // unauthorized
	}

	return s.repo.DeleteComment(uint(idInt))
}