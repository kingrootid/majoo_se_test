package main

import (
	"fmt"
	"log"
	"rootwritter/majoo_test_2_api/internal/database"
	"rootwritter/majoo_test_2_api/internal/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system defaults")
	}

	// Connect to database
	db := database.InitDB()

	// Seed the database with sample data
	err := seedDatabase(db)
	if err != nil {
		log.Fatal("Error seeding database: ", err)
	}

	fmt.Println("Database seeded successfully!")
}

func seedDatabase(db *gorm.DB) error {
	// Check if users already exist to prevent duplicates
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	
	if userCount > 0 {
		fmt.Println("Database already seeded, skipping...")
		return nil
	}

	// Hash passwords
	hashedPassword1, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("mypassword"), bcrypt.DefaultCost)
	hashedPassword3, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	// Create sample users
	users := []models.User{
		{
			Username: "johndoe",
			Email:    "john@example.com",
			Password: string(hashedPassword1),
		},
		{
			Username: "janedoe",
			Email:    "jane@example.com",
			Password: string(hashedPassword2),
		},
		{
			Username: "admin",
			Email:    "admin@example.com",
			Password: string(hashedPassword3),
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("error creating user %s: %v", user.Username, err)
		}
		fmt.Printf("Created user: %s\n", user.Username)
	}

	// Get users from DB to use their IDs
	var createdUsers []models.User
	db.Find(&createdUsers)

	// Create sample posts
	posts := []models.Post{
		{
			Title:   "Welcome to My Blog",
			Content: "This is the first post on my new blog. Excited to share my thoughts with the world!",
			UserID:  createdUsers[0].ID,
		},
		{
			Title:   "Learning Go Programming",
			Content: "Go is an amazing language for building scalable applications. Today I learned about goroutines and channels.",
			UserID:  createdUsers[0].ID,
		},
		{
			Title:   "Travel Adventures",
			Content: "Visited the beautiful mountains last weekend. The view was breathtaking and the air was fresh.",
			UserID:  createdUsers[1].ID,
		},
		{
			Title:   "Admin Announcements",
			Content: "This is an important announcement from the admin. Please read carefully and take note of the updates.",
			UserID:  createdUsers[2].ID,
		},
	}

	for _, post := range posts {
		if err := db.Create(&post).Error; err != nil {
			return fmt.Errorf("error creating post '%s': %v", post.Title, err)
		}
		fmt.Printf("Created post: %s\n", post.Title)
	}

	// Get posts from DB to use their IDs
	var createdPosts []models.Post
	db.Preload("User").Find(&createdPosts)

	// Create sample comments
	comments := []models.Comment{
		{
			Content: "Great post! Looking forward to more content like this.",
			PostID:  createdPosts[0].ID,
			UserID:  createdUsers[1].ID,
		},
		{
			Content: "Thanks for sharing. Very informative!",
			PostID:  createdPosts[1].ID,
			UserID:  createdUsers[2].ID,
		},
		{
			Content: "Amazing experience! Would love to visit there too.",
			PostID:  createdPosts[2].ID,
			UserID:  createdUsers[0].ID,
		},
		{
			Content: "Thank you for the update. Will review and get back with feedback.",
			PostID:  createdPosts[3].ID,
			UserID:  createdUsers[0].ID,
		},
	}

	for _, comment := range comments {
		if err := db.Create(&comment).Error; err != nil {
			return fmt.Errorf("error creating comment: %v", err)
		}
		fmt.Printf("Created comment on post: %s\n", createdPosts[findPostIndex(createdPosts, comment.PostID)].Title)
	}

	return nil
}

// Helper function to find the index of a post by ID
func findPostIndex(posts []models.Post, postID uint) int {
	for i, post := range posts {
		if post.ID == postID {
			return i
		}
	}
	return 0 // fallback
}