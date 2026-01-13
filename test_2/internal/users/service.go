package users

import (
	"errors"
	"strconv"

	"rootwritter/majoo_test_2_api/internal/models"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// CustomValidator implements custom validation rules
type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	// Register custom validation rules if needed
	return &CustomValidator{validator: v}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// validateUserInput performs custom validation for user input
func validateUserInput(username, email, password string) error {
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	// Basic email validation could be added here

	return nil
}

type Service interface {
	RegisterUser(username, email, password string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	UpdateUser(id string, username, email *string, userID uint) (*models.User, error)
	DeleteUser(id string, userID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) RegisterUser(username, email, password string) (*models.User, error) {
	// Validate input
	if err := validateUserInput(username, email, password); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	err = s.repo.CreateUser(user)
	return user, err
}

func (s *service) GetUserByID(id string) (*models.User, error) {
	// Convert id string to uint
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	// Don't return the password hash
	user.Password = ""
	return user, nil
}

func (s *service) UpdateUser(id string, username, email *string, userID uint) (*models.User, error) {
	// Convert id string to uint
	userIDToUpdate, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, err
	}

	// Check if the user is authorized to update this user (can only update own profile)
	if uint(userIDToUpdate) != userID {
		return nil, errors.New("unauthorized: can only update own profile")
	}

	updateData := make(map[string]interface{})
	if username != nil {
		updateData["username"] = *username
	}
	if email != nil {
		updateData["email"] = *email
	}

	updatedUser, err := s.repo.UpdateUser(uint(userIDToUpdate), updateData)
	if err != nil {
		return nil, err
	}

	// Don't return the password hash
	updatedUser.Password = ""
	return updatedUser, nil
}

func (s *service) DeleteUser(id string, userID uint) error {
	// Convert id string to uint
	userIDToDelete, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}

	// Check if the user is authorized to delete this user (can only delete own account)
	if uint(userIDToDelete) != userID {
		return errors.New("unauthorized: can only delete own account")
	}

	return s.repo.DeleteUser(uint(userIDToDelete))
}