package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"lab04-backend/models"

	"github.com/georgysavva/scany/v2/sqlscan"
)

// PostRepository handles database operations for posts
// This repository demonstrates SCANY MAPPING approach for result scanning
type PostRepository struct {
	db *sql.DB
}

// NewPostRepository creates a new PostRepository
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create method using scany for result mapping
func (r *PostRepository) Create(req *models.CreatePostRequest) (*models.Post, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %v", err)
	}

	// Insert into posts table with RETURNING clause
	query := `
		INSERT INTO posts (user_id, title, content, published, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, user_id, title, content, published, created_at, updated_at
	`

	var post models.Post
	err := sqlscan.Get(context.Background(), r.db, &post, query,
		req.UserID, req.Title, req.Content, req.Published)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %v", err)
	}

	return &post, nil
}

// GetByID method using scany
func (r *PostRepository) GetByID(id int) (*models.Post, error) {
	query := `
		SELECT id, user_id, title, content, published, created_at, updated_at
		FROM posts WHERE id = $1
	`

	var post models.Post
	err := sqlscan.Get(context.Background(), r.db, &post, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get post: %v", err)
	}

	return &post, nil
}

// GetByUserID method using scany
func (r *PostRepository) GetByUserID(userID int) ([]models.Post, error) {
	query := `
		SELECT id, user_id, title, content, published, created_at, updated_at
		FROM posts WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var posts []models.Post
	err := sqlscan.Select(context.Background(), r.db, &posts, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by user: %v", err)
	}

	return posts, nil
}

// GetPublished method using scany
func (r *PostRepository) GetPublished() ([]models.Post, error) {
	query := `
		SELECT id, user_id, title, content, published, created_at, updated_at
		FROM posts WHERE published = true
		ORDER BY created_at DESC
	`

	var posts []models.Post
	err := sqlscan.Select(context.Background(), r.db, &posts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get published posts: %v", err)
	}

	return posts, nil
}

// GetAll method using scany
func (r *PostRepository) GetAll() ([]models.Post, error) {
	query := `
		SELECT id, user_id, title, content, published, created_at, updated_at
		FROM posts ORDER BY created_at DESC
	`

	var posts []models.Post
	err := sqlscan.Select(context.Background(), r.db, &posts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all posts: %v", err)
	}

	return posts, nil
}

// Update method using scany
func (r *PostRepository) Update(id int, req *models.UpdatePostRequest) (*models.Post, error) {
	// Build dynamic UPDATE query based on non-nil fields
	var setParts []string
	var args []interface{}
	argIndex := 1

	if req.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *req.Title)
		argIndex++
	}

	if req.Content != nil {
		setParts = append(setParts, fmt.Sprintf("content = $%d", argIndex))
		args = append(args, *req.Content)
		argIndex++
	}

	if req.Published != nil {
		setParts = append(setParts, fmt.Sprintf("published = $%d", argIndex))
		args = append(args, *req.Published)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Add updated_at and id to args
	setParts = append(setParts, fmt.Sprintf("updated_at = NOW()"))
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE posts SET %s
		WHERE id = $%d
		RETURNING id, user_id, title, content, published, created_at, updated_at
	`, strings.Join(setParts, ", "), argIndex)

	var post models.Post
	err := sqlscan.Get(context.Background(), r.db, &post, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %v", err)
	}

	return &post, nil
}

// Delete method (standard SQL)
func (r *PostRepository) Delete(id int) error {
	query := `DELETE FROM posts WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}

// Count method (standard SQL)
func (r *PostRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM posts`

	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count posts: %v", err)
	}

	return count, nil
}

// CountByUserID method (standard SQL)
func (r *PostRepository) CountByUserID(userID int) (int, error) {
	query := `SELECT COUNT(*) FROM posts WHERE user_id = $1`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count posts by user: %v", err)
	}

	return count, nil
}
