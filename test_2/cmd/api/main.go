package main

import (
	"log"
	"net/http"

	"rootwritter/majoo_test_2_api/internal/database"
	"rootwritter/majoo_test_2_api/internal/responses"
	"rootwritter/majoo_test_2_api/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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
		c.JSON(http.StatusNotFound, responses.ErrorResponse{
			Error: "Endpoint not found",
			Code:  http.StatusNotFound,
		})
	})

	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, responses.ErrorResponse{
			Error: "Method not allowed",
			Code:  http.StatusMethodNotAllowed,
		})
	})

	// 5. Health Check Endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, responses.NewSuccess("pong", nil))
	})

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
