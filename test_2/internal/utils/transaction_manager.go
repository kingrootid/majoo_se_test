package utils

import (
	"rootwritter/majoo_test_2_api/internal/comments"
	"rootwritter/majoo_test_2_api/internal/posts"
	"rootwritter/majoo_test_2_api/internal/users"

	"gorm.io/gorm"
)

// TransactionManager manages database transactions across repositories
type TransactionManager struct {
	userRepo    users.Repository
	postRepo    posts.Repository
	commentRepo comments.Repository
	db          *gorm.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(
	userRepo users.Repository,
	postRepo posts.Repository,
	commentRepo comments.Repository,
	db *gorm.DB,
) *TransactionManager {
	return &TransactionManager{
		userRepo:    userRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
		db:          db,
	}
}

// WithTransaction executes the provided function within a database transaction
func (tm *TransactionManager) WithTransaction(fn func(*TransactionManager) error) error {
	return tm.db.Transaction(func(tx *gorm.DB) error {
		// Create new repositories with the transaction DB
		txUserRepo := tm.userRepo.WithTransaction(tx)
		txPostRepo := tm.postRepo.WithTransaction(tx)
		txCommentRepo := tm.commentRepo.WithTransaction(tx)

		// Create a new transaction manager with transaction repositories
		txManager := &TransactionManager{
			userRepo:    txUserRepo,
			postRepo:    txPostRepo,
			commentRepo: txCommentRepo,
			db:          tx,
		}

		return fn(txManager)
	})
}

// GetUserRepo returns the user repository
func (tm *TransactionManager) GetUserRepo() users.Repository {
	return tm.userRepo
}

// GetPostRepo returns the post repository
func (tm *TransactionManager) GetPostRepo() posts.Repository {
	return tm.postRepo
}

// GetCommentRepo returns the comment repository
func (tm *TransactionManager) GetCommentRepo() comments.Repository {
	return tm.commentRepo
}