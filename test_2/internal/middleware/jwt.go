package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"rootwritter/majoo_test_2_api/internal/responses"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"rootwritter/majoo_test_2_api/internal/models"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Claims represents the JWT claims
type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

// JWTMiddleware validates the JWT token
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, responses.NewError("Authorization header required", 401))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(401, responses.NewError("Bearer token required", 401))
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, responses.NewError("Invalid token", 401))
			c.Abort()
			return
		}

		// Additional validation for claims
		if claims.UserID == 0 {
			c.JSON(401, responses.NewError("Invalid token claims", 401))
			c.Abort()
			return
		}

		// Store user ID in context for use in handlers
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// @Summary User login
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login Credentials"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ValidationErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /login [post]
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler handles user login
func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			// Handle gin validation errors
			validationErrors := []responses.ValidationErrorDetail{}

			// Check if it's a validator.ValidationErrors type
			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				for _, fieldErr := range validationErrs {
					validationErrors = append(validationErrors, responses.ValidationErrorDetail{
						Field:   fieldErr.Field(),
						Message: getFieldValidationErrorMessage(fieldErr),
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

		// Find user in database
		var user models.User
		result := db.Where("username = ?", credentials.Username).First(&user)
		if result.Error != nil {
			// Log the error for debugging (remove in production)
			fmt.Printf("Error finding user: %v\n", result.Error)
			c.JSON(401, responses.NewError("Invalid credentials", 401))
			return
		}

		// Validate password using bcrypt
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
			// Log the error for debugging (remove in production)
			fmt.Printf("Error comparing password: %v\n", err)
			c.JSON(401, responses.NewError("Invalid credentials", 401))
			return
		}

		// Generate token
		token, err := GenerateToken(user.ID)
		if err != nil {
			c.JSON(500, responses.NewError("Could not generate token", 500))
			return
		}

		response := map[string]interface{}{
			"token": token,
			"user":  user.Username,
		}

		c.JSON(200, responses.NewSuccess("Login successful", response))
	}
}

// getFieldValidationErrorMessage returns a human-readable validation error message
func getFieldValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	default:
		return fe.Tag()
	}
}