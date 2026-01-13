package models

import (
	"time"

	"gorm.io/gorm"
)

// User Table
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"unique;not null" json:"username"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"` // "-" agar password tidak muncul di JSON
	Posts     []Post         `json:"posts"`             // Relasi One-to-Many ke Post
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Post Table
type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	UserID    uint           `json:"user_id"`  // Foreign Key ke User
	Comments  []Comment      `json:"comments"` // Relasi One-to-Many ke Comment
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Comment Table
type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Content   string         `gorm:"not null" json:"content"`
	PostID    uint           `json:"post_id"` // Foreign Key ke Post
	UserID    uint           `json:"user_id"` // Foreign Key ke User (siapa yang komen)
	Post      Post           `json:"post" gorm:"foreignKey:PostID"`  // Relasi ke Post
	User      User           `json:"user" gorm:"foreignKey:UserID"`  // Relasi ke User
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
