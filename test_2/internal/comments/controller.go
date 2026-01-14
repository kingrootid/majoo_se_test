package comments

import (
	"fmt"
	"net/http"

	"rootwritter/majoo_test_2_api/internal/responses"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// @Summary Create a new comment
// @Description Create a new comment for a specific post
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param comment body CreateCommentRequest true "Comment Data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ValidationErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /posts/{id}/comments [post]
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// @Summary Get comments by post ID
// @Description Get all comments for a specific post with pagination
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param page query string false "Page number" default 1
// @Param limit query string false "Number of items per page" default 10
// @Success 200 {object} responses.SuccessResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /posts/{id}/comments [get]

// @Summary Update a comment by ID
// @Description Update a specific comment by its ID if the user owns it
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Param comment body UpdateCommentRequest true "Update Comment Data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ValidationErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /comments/{id} [put]
type UpdateCommentRequest struct {
	Content *string `json:"content" binding:"omitempty,min=1,max=1000"`
}

// @Summary Delete a comment by ID
// @Description Delete a specific comment by its ID if the user owns it
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /comments/{id} [delete]

type Controller struct {
	svc Service
}

func NewController(svc Service) *Controller {
	return &Controller{svc}
}

// @Summary Create a new comment
// @Description Create a new comment for a specific post
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param comment body CreateCommentRequest true "Comment Data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ValidationErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /posts/{id}/comments [post]
func (ctrl *Controller) Create(c *gin.Context) {
	// Get the post ID from the URL parameter (from nested route structure)
	postIDStr := c.Param("id")
	var postID uint
	_, err := fmt.Sscan(postIDStr, &postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewError("Invalid post ID", http.StatusBadRequest))
		return
	}

	var input struct {
		Content string `json:"content" binding:"required,min=1,max=1000"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		// Handle gin validation errors
		validationErrors := []responses.ValidationErrorDetail{}
		
		// Check if it's a validator.ValidationErrors type
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrs {
				validationErrors = append(validationErrors, responses.ValidationErrorDetail{
					Field:   fieldErr.Field(),
					Message: getCommentValidationErrorMessage(fieldErr),
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

	// Get user ID from JWT middleware
	userID := c.MustGet("userID").(uint)

	comment, err := ctrl.svc.CreateNewComment(input.Content, postID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewError("Failed to create comment", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, responses.NewSuccess("Comment created successfully", comment))
}

// getCommentValidationErrorMessage returns a human-readable validation error message for comments
func getCommentValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		min := fe.Param()
		return fmt.Sprintf("Value must be at least %s characters long", min)
	case "max":
		max := fe.Param()
		return fmt.Sprintf("Value must be at most %s characters long", max)
	default:
		return fe.Tag()
	}
}

// @Summary Get a comment by ID
// @Description Get a specific comment by its ID
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /comments/{id} [get]
func (ctrl *Controller) GetByID(c *gin.Context) {
	id := c.Param("id")

	comment, err := ctrl.svc.GetCommentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.NewError("Comment not found", http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccess("Comment retrieved successfully", comment))
}

// @Summary Get comments by post ID
// @Description Get all comments for a specific post with pagination
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param page query string false "Page number" default 1
// @Param limit query string false "Number of items per page" default 10
// @Success 200 {object} responses.SuccessResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /posts/{id}/comments [get]
func (ctrl *Controller) GetByPostID(c *gin.Context) {
	// Get the post ID from the URL parameter (from nested route structure)
	postID := c.Param("id")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	comments, err := ctrl.svc.GetCommentsByPostID(postID, page, limit)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.NewError("Post not found or error retrieving comments", http.StatusNotFound))
		return
	}

	meta := map[string]interface{}{
		"page":  page,
		"limit": limit,
		"total": len(comments),
	}

	c.JSON(http.StatusOK, responses.NewSuccessWithMeta("Comments retrieved successfully", comments, meta))
}

// @Summary Get comments by user ID
// @Description Get all comments by a specific user with pagination
// @Tags comments
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param page query string false "Page number" default 1
// @Param limit query string false "Number of items per page" default 10
// @Success 200 {object} responses.SuccessResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /users/{user_id}/comments [get]
func (ctrl *Controller) GetByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	comments, err := ctrl.svc.GetCommentsByUserID(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.NewError("User not found or error retrieving comments", http.StatusNotFound))
		return
	}

	meta := map[string]interface{}{
		"page":  page,
		"limit": limit,
		"total": len(comments),
	}

	c.JSON(http.StatusOK, responses.NewSuccessWithMeta("Comments retrieved successfully", comments, meta))
}

// @Summary Update a comment by ID
// @Description Update a specific comment by its ID if the user owns it
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Param comment body UpdateCommentRequest true "Update Comment Data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ValidationErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /comments/{id} [put]
func (ctrl *Controller) Update(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Content *string `json:"content" binding:"omitempty,min=1,max=1000"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		// Handle gin validation errors
		validationErrors := []responses.ValidationErrorDetail{}
		
		// Check if it's a validator.ValidationErrors type
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrs {
				validationErrors = append(validationErrors, responses.ValidationErrorDetail{
					Field:   fieldErr.Field(),
					Message: getCommentValidationErrorMessage(fieldErr),
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

	// Get user ID from JWT middleware
	userID := c.MustGet("userID").(uint)

	updatedComment, err := ctrl.svc.UpdateComment(id, input.Content, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusUnauthorized, responses.NewError("Not authorized to update this comment", http.StatusUnauthorized))
			return
		}
		if err.Error() == "content cannot be empty" {
			c.JSON(http.StatusBadRequest, responses.NewError("Content cannot be empty", http.StatusBadRequest))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewError("Failed to update comment", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccess("Comment updated successfully", updatedComment))
}

// @Summary Delete a comment by ID
// @Description Delete a specific comment by its ID if the user owns it
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /comments/{id} [delete]
func (ctrl *Controller) Delete(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from JWT middleware
	userID := c.MustGet("userID").(uint)

	err := ctrl.svc.DeleteComment(id, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusUnauthorized, responses.NewError("Not authorized to delete this comment", http.StatusUnauthorized))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewError("Failed to delete comment", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccess("Comment deleted successfully", nil))
}