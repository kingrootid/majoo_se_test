package routes

import (
	"rootwritter/majoo_test_2_api/internal/comments"
	"rootwritter/majoo_test_2_api/internal/middleware"
	"rootwritter/majoo_test_2_api/internal/posts"
	"rootwritter/majoo_test_2_api/internal/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// 1. Initialize Layers
	userRepo := users.NewRepository(db)
	userSvc := users.NewService(userRepo)
	userCtrl := users.NewController(userSvc)

	postRepo := posts.NewRepository(db)
	postSvc := posts.NewService(postRepo)
	postCtrl := posts.NewController(postSvc)

	commentRepo := comments.NewRepository(db)
	commentSvc := comments.NewService(commentRepo)
	commentCtrl := comments.NewController(commentSvc)

	// 2. Define Routes
	api := r.Group("/api/v1")
	{
		// Public routes
		api.POST("/register", userCtrl.Register)
		api.POST("/login", middleware.LoginHandler(db))

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.JWTMiddleware())
		{
			// User routes
			protected.GET("/profile", userCtrl.GetProfile)
			protected.PUT("/profile", userCtrl.UpdateProfile)
			protected.DELETE("/profile", userCtrl.DeleteAccount)

			// Posts routes need to come first with their sub-routes before individual post routes
			postsGroup := protected.Group("/posts")
			{
				postsGroup.POST("", postCtrl.Create)                                 // Create a new post
				postsGroup.GET("", postCtrl.GetAll)                                  // Get all posts

				// Nested routes for post-specific operations (this avoids conflicts)
				singlePostGroup := postsGroup.Group("/:id")
				{
					singlePostGroup.GET("", postCtrl.GetByID)                         // Get a specific post
					singlePostGroup.PUT("", postCtrl.Update)                          // Update a specific post
					singlePostGroup.DELETE("", postCtrl.Delete)                       // Delete a specific post

					// Comments related to a specific post
					singlePostGroup.POST("/comments", commentCtrl.Create)             // Create comment on a post
					singlePostGroup.GET("/comments", commentCtrl.GetByPostID)         // Get all comments for a post
				}
			}

			// Individual comment routes
			commentsGroup := protected.Group("/comments")
			{
				commentsGroup.GET("/:id", commentCtrl.GetByID)                       // Get specific comment
				commentsGroup.PUT("/:id", commentCtrl.Update)                        // Update a comment
				commentsGroup.DELETE("/:id", commentCtrl.Delete)                     // Delete a comment
			}

			// Get comments by user
			protected.GET("/users/:user_id/comments", commentCtrl.GetByUserID)       // Get all comments by a user
		}
	}
}
