// @title           Rest API Development (Majoo SE Test)
// @version         1.0
// @description     This is a sample API with Gin and Swagger.
// @host              localhost:8090
// @BasePath         /api/v1
package main

import (
	"log"
	"net/http"

	"rootwritter/majoo_test_2_api/internal/database"
	"rootwritter/majoo_test_2_api/internal/responses"
	"rootwritter/majoo_test_2_api/internal/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/joho/godotenv"
	_ "rootwritter/majoo_test_2_api/docs" // Import for Swagger documentation
)

// @title Blog API Documentation
// @version 1.0
// @description Complete REST API for a blog system with user authentication, posts, and comments
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token

// @host localhost:8090
// @BasePath /api/v1
// @schemes http
func main() {
	// 1. Load Environment Variables (Config)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system default")
	}

	// 2. Koneksi ke Database & Auto-Migration
	db := database.InitDB()

	// 3. Inisialisasi Framework (Gin)
	r := gin.Default()

	// 4. Global Middleware (CORS, Recovery, dsb)
	r.Use(gin.Recovery())

	// Global error handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, responses.NewError("Endpoint not found", http.StatusNotFound))
	})

	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, responses.NewError("Method not allowed", http.StatusMethodNotAllowed))
	})

	// 5. Health Check Endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, responses.NewSuccess("pong", nil))
	})

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 6. Registrasi Group Routes (internal/routes)
	// Kita oper 'db' agar handler bisa mengakses database
	routes.SetupRoutes(r, db)

	// 7. Jalankan Server
	port := ":8090"
	log.Printf("Server running on http://localhost%s", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
