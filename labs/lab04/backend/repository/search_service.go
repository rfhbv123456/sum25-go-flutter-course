package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"lab04-backend/models"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/sqlscan"
)

// SearchService handles dynamic search operations using Squirrel query builder
// This service demonstrates SQUIRREL QUERY BUILDER approach for dynamic SQL
type SearchService struct {
	db   *sql.DB
	psql squirrel.StatementBuilderType
}

// SearchFilters represents search parameters
type SearchFilters struct {
	Query        string // Search in title and content
	UserID       *int   // Filter by user ID
	Published    *bool  // Filter by published status
	MinWordCount *int   // Minimum word count in content
	Limit        int    // Results limit (default 50)
	Offset       int    // Results offset (for pagination)
	OrderBy      string // Order by field (title, created_at, updated_at)
	OrderDir     string // Order direction (ASC, DESC)
}

// NewSearchService creates a new SearchService
func NewSearchService(db *sql.DB) *SearchService {
	return &SearchService{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// SearchPosts method using Squirrel query builder
func (s *SearchService) SearchPosts(ctx context.Context, filters SearchFilters) ([]models.Post, error) {
	// Start with base query
	query := s.psql.Select("id", "user_id", "title", "content", "published", "created_at", "updated_at").
		From("posts")

	// Add WHERE conditions dynamically
	if filters.Query != "" {
		searchTerm := "%" + filters.Query + "%"
		query = query.Where(squirrel.Or{
			squirrel.ILike{"title": searchTerm},
			squirrel.ILike{"content": searchTerm},
		})
	}

	if filters.UserID != nil {
		query = query.Where(squirrel.Eq{"user_id": *filters.UserID})
	}

	if filters.Published != nil {
		query = query.Where(squirrel.Eq{"published": *filters.Published})
	}

	if filters.MinWordCount != nil {
		// Count words in content (basic implementation)
		wordCountExpr := fmt.Sprintf("array_length(string_to_array(content, ' '), 1) >= %d", *filters.MinWordCount)
		query = query.Where(wordCountExpr)
	}

	// Add ORDER BY dynamically
	if filters.OrderBy != "" {
		validOrderFields := map[string]bool{
			"title":      true,
			"created_at": true,
			"updated_at": true,
			"user_id":    true,
		}

		if validOrderFields[filters.OrderBy] {
			orderDir := "DESC"
			if strings.ToUpper(filters.OrderDir) == "ASC" {
				orderDir = "ASC"
			}
			query = query.OrderBy(filters.OrderBy + " " + orderDir)
		} else {
			query = query.OrderBy("created_at DESC")
		}
	} else {
		query = query.OrderBy("created_at DESC")
	}

	// Add LIMIT/OFFSET
	if filters.Limit > 0 {
		query = query.Limit(uint64(filters.Limit))
	} else {
		query = query.Limit(50) // Default limit
	}

	if filters.Offset > 0 {
		query = query.Offset(uint64(filters.Offset))
	}

	// Build final SQL
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	// Execute with scany
	var posts []models.Post
	err = sqlscan.Select(ctx, s.db, &posts, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %v", err)
	}

	return posts, nil
}

// SearchUsers method using Squirrel
func (s *SearchService) SearchUsers(ctx context.Context, nameQuery string, limit int) ([]models.User, error) {
	query := s.psql.Select("id", "name", "email", "created_at", "updated_at").
		From("users").
		Where(squirrel.Like{"name": "%" + nameQuery + "%"}).
		OrderBy("name").
		Limit(uint64(limit))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build user search query: %v", err)
	}

	var users []models.User
	err = sqlscan.Select(ctx, s.db, &users, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %v", err)
	}

	return users, nil
}

// GetPostStats method using Squirrel with JOINs
func (s *SearchService) GetPostStats(ctx context.Context) (*PostStats, error) {
	query := s.psql.Select(
		"COUNT(p.id) as total_posts",
		"COUNT(CASE WHEN p.published = true THEN 1 END) as published_posts",
		"COUNT(DISTINCT p.user_id) as active_users",
		"AVG(LENGTH(p.content)) as avg_content_length",
	).From("posts p").
		Join("users u ON p.user_id = u.id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build stats query: %v", err)
	}

	var stats PostStats
	err = sqlscan.Get(ctx, s.db, &stats, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get post stats: %v", err)
	}

	return &stats, nil
}

// PostStats represents aggregated post statistics
type PostStats struct {
	TotalPosts       int     `db:"total_posts"`
	PublishedPosts   int     `db:"published_posts"`
	ActiveUsers      int     `db:"active_users"`
	AvgContentLength float64 `db:"avg_content_length"`
}

// BuildDynamicQuery helper method
func (s *SearchService) BuildDynamicQuery(baseQuery squirrel.SelectBuilder, filters SearchFilters) squirrel.SelectBuilder {
	query := baseQuery

	if filters.Query != "" {
		searchTerm := "%" + filters.Query + "%"
		query = query.Where(squirrel.Or{
			squirrel.ILike{"title": searchTerm},
			squirrel.ILike{"content": searchTerm},
		})
	}

	if filters.UserID != nil {
		query = query.Where(squirrel.Eq{"user_id": *filters.UserID})
	}

	if filters.Published != nil {
		query = query.Where(squirrel.Eq{"published": *filters.Published})
	}

	if filters.MinWordCount != nil {
		wordCountExpr := fmt.Sprintf("array_length(string_to_array(content, ' '), 1) >= %d", *filters.MinWordCount)
		query = query.Where(wordCountExpr)
	}

	// Add ORDER BY
	if filters.OrderBy != "" {
		validOrderFields := map[string]bool{
			"title":      true,
			"created_at": true,
			"updated_at": true,
			"user_id":    true,
		}

		if validOrderFields[filters.OrderBy] {
			orderDir := "DESC"
			if strings.ToUpper(filters.OrderDir) == "ASC" {
				orderDir = "ASC"
			}
			query = query.OrderBy(filters.OrderBy + " " + orderDir)
		} else {
			query = query.OrderBy("created_at DESC")
		}
	} else {
		query = query.OrderBy("created_at DESC")
	}

	// Add LIMIT/OFFSET
	if filters.Limit > 0 {
		query = query.Limit(uint64(filters.Limit))
	} else {
		query = query.Limit(50)
	}

	if filters.Offset > 0 {
		query = query.Offset(uint64(filters.Offset))
	}

	return query
}

// GetTopUsers method using Squirrel with complex aggregation
func (s *SearchService) GetTopUsers(ctx context.Context, limit int) ([]UserWithStats, error) {
	query := s.psql.Select(
		"u.id",
		"u.name",
		"u.email",
		"u.created_at",
		"u.updated_at",
		"COUNT(p.id) as post_count",
		"COUNT(CASE WHEN p.published = true THEN 1 END) as published_count",
		"MAX(p.created_at) as last_post_date",
	).From("users u").
		LeftJoin("posts p ON u.id = p.user_id").
		GroupBy("u.id", "u.name", "u.email", "u.created_at", "u.updated_at").
		OrderBy("post_count DESC").
		Limit(uint64(limit))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build top users query: %v", err)
	}

	var users []UserWithStats
	err = sqlscan.Select(ctx, s.db, &users, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get top users: %v", err)
	}

	return users, nil
}

// UserWithStats represents a user with post statistics
type UserWithStats struct {
	models.User
	PostCount      int    `db:"post_count"`
	PublishedCount int    `db:"published_count"`
	LastPostDate   string `db:"last_post_date"`
}
