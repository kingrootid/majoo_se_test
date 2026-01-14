package posts

import (
	"net/http"

	"rootwritter/majoo_test_2_api/internal/responses"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// @Summary Create a new post
// @Description Create a new blog post for authenticated user
// @Tags posts
// @Accept json
// @Produce json
// @Param post body CreatePostRequest true "Post Data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ValidationErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /posts [post]
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// @Summary Update a post by ID
// @Description Update a specific post by its ID if the user owns it
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param post body UpdatePostRequest true "Update Post Data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ValidationErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /posts/{id} [put]
type UpdatePostRequest struct {
	Title   *string `json:"title" binding:"omitempty,min=1,max=100"`
	Content *string `json:"content" binding:"omitempty,min=1,max=10000"`
}

type Controller struct {
	svc Service
}

func NewController(svc Service) *Controller {
	return &Controller{svc}
}

// @Summary Create a new post
// @Description Create a new blog post for authenticated user
// @Tags posts
// @Accept json
// @Produce json
// @Param post body CreatePostRequest true "Post Data"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ValidationErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /posts [post]
func (ctrl *Controller) Create(c *gin.Context) {
	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		// Handle gin validation errors
		validationErrors := []responses.ValidationErrorDetail{}
		
		// Check if it's a validator.ValidationErrors type
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrs {
				validationErrors = append(validationErrors, responses.ValidationErrorDetail{
					Field:   fieldErr.Field(),
					Message: getValidationErrorMessage(fieldErr),
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

	// Ambil ID dari middleware JWT
	userID := c.MustGet("userID").(uint)

	post, err := ctrl.svc.CreateNewPost(input.Title, input.Content, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewError("Failed to create post", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, responses.NewSuccess("Post created successfully", post))
}

// getValidationErrorMessage returns a human-readable validation error message
func getValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Value too short"
	case "max":
		return "Value too long"
	case "email":
		return "Invalid email format"
	case "numeric":
		return "Must be a number"
	default:
		return fe.Tag()
	}
}

// @Summary Get all posts
// @Description Get a list of all posts with pagination
// @Tags posts
// @Accept json
// @Produce json
// @Param page query string false "Page number" default 1
// @Param limit query string false "Number of items per page" default 10
// @Success 200 {object} responses.SuccessResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /posts [get]
func (ctrl *Controller) GetAll(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	posts, err := ctrl.svc.GetAllPosts(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewError("Failed to fetch posts", http.StatusInternalServerError))
		return
	}

	meta := map[string]interface{}{
		"page":  page,
		"limit": limit,
		"total": len(posts),
	}

	c.JSON(http.StatusOK, responses.NewSuccessWithMeta("Posts retrieved successfully", posts, meta))
}

// @Summary Get a post by ID
// @Description Get a specific post by its ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /posts/{id} [get]
func (ctrl *Controller) GetByID(c *gin.Context) {
	id := c.Param("id")

	post, err := ctrl.svc.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.NewError("Post not found", http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccess("Post retrieved successfully", post))
}

// @Summary Update a post by ID
// @Description Update a specific post by its ID if the user owns it
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param post body UpdatePostRequest true "Update Post Data"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ValidationErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /posts/{id} [put]
func (ctrl *Controller) Update(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Title   *string `json:"title" binding:"omitempty,min=1,max=100"`
		Content *string `json:"content" binding:"omitempty,min=1,max=10000"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		// Handle gin validation errors
		validationErrors := []responses.ValidationErrorDetail{}
		
		// Check if it's a validator.ValidationErrors type
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrs {
				validationErrors = append(validationErrors, responses.ValidationErrorDetail{
					Field:   fieldErr.Field(),
					Message: getValidationErrorMessage(fieldErr),
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

	// Ambil ID dari middleware JWT
	userID := c.MustGet("userID").(uint)

	updatedPost, err := ctrl.svc.UpdatePost(id, input.Title, input.Content, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusUnauthorized, responses.NewError("Not authorized to update this post", http.StatusUnauthorized))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewError("Failed to update post", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccess("Post updated successfully", updatedPost))
}

// @Summary Delete a post by ID
// @Description Delete a specific post by its ID if the user owns it
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /posts/{id} [delete]
func (ctrl *Controller) Delete(c *gin.Context) {
	id := c.Param("id")

	// Ambil ID dari middleware JWT
	userID := c.MustGet("userID").(uint)

	err := ctrl.svc.DeletePost(id, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusUnauthorized, responses.NewError("Not authorized to delete this post", http.StatusUnauthorized))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewError("Failed to delete post", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, responses.NewSuccess("Post deleted successfully", nil))
}