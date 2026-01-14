package users

import (
	"fmt"
	"net/http"
	"strconv"

	"rootwritter/majoo_test_2_api/internal/responses"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// @Summary Register a new user
// @Description Register a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body RegisterUserRequest true "User Registration Data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ValidationErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /register [post]
type RegisterUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

type Controller struct {
	svc Service
}

func NewController(svc Service) *Controller {
	return &Controller{svc}
}

// @Summary Register a new user
// @Description Register a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body RegisterUserRequest true "User Registration Data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ValidationErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /register [post]
func (ctrl *Controller) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=100"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		// Handle gin validation errors
		validationErrors := []responses.ValidationErrorDetail{}

		// Check if it's a validator.ValidationErrors type
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrs {
				validationErrors = append(validationErrors, responses.ValidationErrorDetail{
					Field:   fieldErr.Field(),
					Message: getUserValidationErrorMessage(fieldErr),
					Value:   fieldErr.Value(),
				})
			}
		} else {
			// General error
			validationErrors = append(validationErrors, responses.ValidationErrorDetail{
				Field:   "unknown",
				Message: err.Error(),
			})
		}

		c.JSON(http.StatusUnprocessableEntity, responses.NewValidationError("Validation failed", validationErrors))
		return
	}

	user, err := ctrl.svc.RegisterUser(input.Username, input.Email, input.Password)
	if err != nil {
		if err.Error() == "username must be at least 3 characters long" ||
		   err.Error() == "password must be at least 6 characters long" {
			c.JSON(http.StatusBadRequest, responses.NewError("Validation error: " + err.Error(), http.StatusBadRequest))
			return
		}
		// Check for duplicate entry errors in different databases
		if containsError(err.Error(), []string{"UNIQUE constraint failed", "duplicate key value violates unique constraint"}) {
			c.JSON(http.StatusConflict, responses.NewError("Username or email already exists", http.StatusConflict))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewError("Failed to register user: " + err.Error(), http.StatusInternalServerError))
		return
	}

	// Don't return the password hash
	user.Password = ""

	c.JSON(http.StatusCreated, responses.NewSuccess("User registered successfully", user))
}

// getUserValidationErrorMessage returns a human-readable validation error message for users
func getUserValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		min := fe.Param()
		return fmt.Sprintf("Value must be at least %s characters long", min)
	case "max":
		max := fe.Param()
		return fmt.Sprintf("Value must be at most %s characters long", max)
	case "email":
		return "Invalid email format"
	default:
		return fe.Tag()
	}
}

func (ctrl *Controller) GetProfile(c *gin.Context) {
	// Get user ID from JWT middleware
	userID := c.MustGet("userID").(uint)

	user, err := ctrl.svc.GetUserByID(strconv.Itoa(int(userID)))
	if err != nil {
		c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Error: "User not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccess("User profile retrieved successfully", user))
}

func (ctrl *Controller) UpdateProfile(c *gin.Context) {
	// Get user ID from JWT middleware
	userID := c.MustGet("userID").(uint)

	var input struct {
		Username *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
		Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		// Handle gin validation errors
		validationErrors := []responses.ValidationErrorDetail{}

		// Check if it's a validator.ValidationErrors type
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrs {
				validationErrors = append(validationErrors, responses.ValidationErrorDetail{
					Field:   fieldErr.Field(),
					Message: getUserValidationErrorMessage(fieldErr),
					Value:   fieldErr.Value(),
				})
			}
		} else {
			// General error
			validationErrors = append(validationErrors, responses.ValidationErrorDetail{
				Field:   "unknown",
				Message: err.Error(),
			})
		}

		c.JSON(http.StatusUnprocessableEntity, responses.NewValidationError("Validation failed", validationErrors))
		return
	}

	updatedUser, err := ctrl.svc.UpdateUser(strconv.Itoa(int(userID)), input.Username, input.Email, userID)
	if err != nil {
		if err.Error() == "unauthorized: can only update own profile" {
			c.JSON(http.StatusUnauthorized, responses.NewError("Not authorized to update this profile", http.StatusUnauthorized))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewError("Failed to update user", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccess("User updated successfully", updatedUser))
}

func (ctrl *Controller) DeleteAccount(c *gin.Context) {
	// Get user ID from JWT middleware
	userID := c.MustGet("userID").(uint)

	err := ctrl.svc.DeleteUser(strconv.Itoa(int(userID)), userID)
	if err != nil {
		if err.Error() == "unauthorized: can only delete own account" {
			c.JSON(http.StatusUnauthorized, responses.ErrorResponse{
				Error: "Not authorized to delete this account",
				Code:  http.StatusUnauthorized,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error: "Failed to delete user",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccess("User deleted successfully", nil))
}

// Helper function to check if error message contains any of the substrings
func containsError(errStr string, substrs []string) bool {
	for _, substr := range substrs {
		for i := 0; i <= len(errStr)-len(substr); i++ {
			if i+len(substr) <= len(errStr) && errStr[i:i+len(substr)] == substr {
				return true
			}
		}
	}
	return false
}