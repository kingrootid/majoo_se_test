package users

import (
	"errors"
	"strconv"

	"rootwritter/majoo_test_2_api/internal/models"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	// Example of a complex transaction operation
	CreateUserWithProfile(username, email, password string) (*models.User, error)

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

// CreateUserWithProfile demonstrates a complex transaction
func (s *service) CreateUserWithProfile(username, email, password string) (*models.User, error) {
	var user *models.User
	var err error

	// Start a transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Create user in the transaction context
		txService := s.WithTransaction(tx)

		// Register the user
		user, err = txService.RegisterUser(username, email, password)
		if err != nil {
			return err // Rolling back the transaction
		}

		// Here we could perform additional operations in the same transaction
		// For example, creating a user profile or initializing settings
		// Example:
		// profile := &models.UserProfile{UserID: user.ID, Status: "active"}
		// if err := tx.Create(&profile).Error; err != nil {
		//     return err
		// }

		return nil // Commit the transaction
	})

	if err != nil {
		return nil, err
	}

	// Don't return the password hash
	user.Password = ""
	return user, nil
}