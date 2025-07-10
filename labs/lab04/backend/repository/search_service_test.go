package repository

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"lab04-backend/database"
	"lab04-backend/models"

	"github.com/Masterminds/squirrel"
)

// TestSearchService tests the Squirrel query builder approach
func TestSearchService(t *testing.T) {
	// Initialize database for testing
	db, err := database.InitDB()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB(db)

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create service instance
	searchService := NewSearchService(db)

	// Test SearchPosts with various filters
	t.Run("SearchPosts with filters", func(t *testing.T) {
		// Insert test data
		userRepo := NewUserRepository(db)
		postRepo := NewPostRepository(db)

		// Create test user
		user, err := userRepo.Create(&models.CreateUserRequest{
			Name:  "Test User",
			Email: "test@example.com",
		})
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// Create test posts
		posts := []*models.CreatePostRequest{
			{
				UserID:    user.ID,
				Title:     "Golang Programming",
				Content:   "This is a post about Golang programming language",
				Published: true,
			},
			{
				UserID:    user.ID,
				Title:     "Flutter Development",
				Content:   "This is a post about Flutter development",
				Published: false,
			},
			{
				UserID:    user.ID,
				Title:     "Database Design",
				Content:   "This is a post about database design with SQLite",
				Published: true,
			},
		}

		for _, postReq := range posts {
			_, err := postRepo.Create(postReq)
			if err != nil {
				t.Fatalf("Failed to create test post: %v", err)
			}
		}

		// Test empty filters (should return all posts)
		filters := SearchFilters{}
		searchPosts, err := searchService.SearchPosts(context.Background(), filters)
		if err != nil {
			t.Errorf("SearchPosts with empty filters failed: %v", err)
		}
		if len(searchPosts) < 3 {
			t.Errorf("Expected at least 3 posts, got %d", len(searchPosts))
		}

		// Test search by query string
		filters = SearchFilters{Query: "golang"}
		searchPosts, err = searchService.SearchPosts(context.Background(), filters)
		if err != nil {
			t.Errorf("SearchPosts with query filter failed: %v", err)
		}
		if len(searchPosts) == 0 {
			t.Error("Expected posts with 'golang' query, got none")
		}

		// Test filter by published status
		published := true
		filters = SearchFilters{Published: &published}
		searchPosts, err = searchService.SearchPosts(context.Background(), filters)
		if err != nil {
			t.Errorf("SearchPosts with published filter failed: %v", err)
		}
		for _, post := range searchPosts {
			if !post.Published {
				t.Error("Expected only published posts")
			}
		}

		// Test pagination
		filters = SearchFilters{Limit: 2}
		searchPosts, err = searchService.SearchPosts(context.Background(), filters)
		if err != nil {
			t.Errorf("SearchPosts with limit failed: %v", err)
		}
		if len(searchPosts) > 2 {
			t.Errorf("Expected at most 2 posts, got %d", len(searchPosts))
		}

		// Test sorting
		filters = SearchFilters{OrderBy: "title", OrderDir: "ASC"}
		searchPosts, err = searchService.SearchPosts(context.Background(), filters)
		if err != nil {
			t.Errorf("SearchPosts with sorting failed: %v", err)
		}
		if len(searchPosts) < 2 {
			t.Error("Expected at least 2 posts for sorting test")
		}
	})

	// Test SearchUsers functionality
	t.Run("SearchUsers", func(t *testing.T) {
		// Insert test users
		userRepo := NewUserRepository(db)
		users := []*models.CreateUserRequest{
			{Name: "John Smith", Email: "john@example.com"},
			{Name: "Jane Doe", Email: "jane@example.com"},
			{Name: "Bob Johnson", Email: "bob@example.com"},
		}

		for _, userReq := range users {
			_, err := userRepo.Create(userReq)
			if err != nil {
				t.Fatalf("Failed to create test user: %v", err)
			}
		}

		// Test exact name matches
		searchUsers, err := searchService.SearchUsers(context.Background(), "John", 10)
		if err != nil {
			t.Errorf("SearchUsers failed: %v", err)
		}
		if len(searchUsers) == 0 {
			t.Error("Expected users with 'John' in name, got none")
		}

		// Test partial name matches
		searchUsers, err = searchService.SearchUsers(context.Background(), "Jo", 10)
		if err != nil {
			t.Errorf("SearchUsers with partial match failed: %v", err)
		}
		if len(searchUsers) == 0 {
			t.Error("Expected users with 'Jo' in name, got none")
		}

		// Test limit functionality
		searchUsers, err = searchService.SearchUsers(context.Background(), "J", 2)
		if err != nil {
			t.Errorf("SearchUsers with limit failed: %v", err)
		}
		if len(searchUsers) > 2 {
			t.Errorf("Expected at most 2 users, got %d", len(searchUsers))
		}
	})

	// Test GetPostStats with complex aggregation
	t.Run("GetPostStats", func(t *testing.T) {
		// Insert test data
		userRepo := NewUserRepository(db)
		postRepo := NewPostRepository(db)

		// Create users
		user1, err := userRepo.Create(&models.CreateUserRequest{
			Name:  "User 1",
			Email: "user1@example.com",
		})
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		user2, err := userRepo.Create(&models.CreateUserRequest{
			Name:  "User 2",
			Email: "user2@example.com",
		})
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// Create posts
		posts := []*models.CreatePostRequest{
			{UserID: user1.ID, Title: "Post 1", Content: "Content 1", Published: true},
			{UserID: user1.ID, Title: "Post 2", Content: "Content 2", Published: true},
			{UserID: user2.ID, Title: "Post 3", Content: "Content 3", Published: false},
		}

		for _, postReq := range posts {
			_, err := postRepo.Create(postReq)
			if err != nil {
				t.Fatalf("Failed to create test post: %v", err)
			}
		}

		// Test aggregation
		stats, err := searchService.GetPostStats(context.Background())
		if err != nil {
			t.Errorf("GetPostStats failed: %v", err)
		}

		if stats.TotalPosts < 3 {
			t.Errorf("Expected at least 3 total posts, got %d", stats.TotalPosts)
		}

		if stats.PublishedPosts < 2 {
			t.Errorf("Expected at least 2 published posts, got %d", stats.PublishedPosts)
		}

		if stats.ActiveUsers < 2 {
			t.Errorf("Expected at least 2 active users, got %d", stats.ActiveUsers)
		}
	})

	// Test GetTopUsers with aggregation and sorting
	t.Run("GetTopUsers", func(t *testing.T) {
		// Insert test data
		userRepo := NewUserRepository(db)
		postRepo := NewPostRepository(db)

		// Create users with different post counts
		user1, err := userRepo.Create(&models.CreateUserRequest{
			Name:  "High Poster",
			Email: "high@example.com",
		})
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		user2, err := userRepo.Create(&models.CreateUserRequest{
			Name:  "Low Poster",
			Email: "low@example.com",
		})
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// Create posts for user1 (more posts)
		for i := 0; i < 3; i++ {
			_, err := postRepo.Create(&models.CreatePostRequest{
				UserID:    user1.ID,
				Title:     fmt.Sprintf("Post %d", i+1),
				Content:   fmt.Sprintf("Content %d", i+1),
				Published: true,
			})
			if err != nil {
				t.Fatalf("Failed to create test post: %v", err)
			}
		}

		// Create posts for user2 (fewer posts)
		for i := 0; i < 1; i++ {
			_, err := postRepo.Create(&models.CreatePostRequest{
				UserID:    user2.ID,
				Title:     fmt.Sprintf("Post %d", i+1),
				Content:   fmt.Sprintf("Content %d", i+1),
				Published: true,
			})
			if err != nil {
				t.Fatalf("Failed to create test post: %v", err)
			}
		}

		// Test user ranking
		topUsers, err := searchService.GetTopUsers(context.Background(), 10)
		if err != nil {
			t.Errorf("GetTopUsers failed: %v", err)
		}

		if len(topUsers) == 0 {
			t.Error("Expected top users, got none")
		}

		// Check that users are ordered by post count (descending)
		if len(topUsers) > 1 && topUsers[0].PostCount < topUsers[1].PostCount {
			t.Error("Expected users to be ordered by post count (descending)")
		}
	})

	// Test BuildDynamicQuery helper
	t.Run("BuildDynamicQuery", func(t *testing.T) {
		// Test with different filter combinations
		baseQuery := searchService.psql.Select("*").From("posts")

		// Test with query filter
		filters := SearchFilters{Query: "test"}
		query := searchService.BuildDynamicQuery(baseQuery, filters)
		sql, _, err := query.ToSql()
		if err != nil {
			t.Errorf("BuildDynamicQuery failed: %v", err)
		}
		if !strings.Contains(sql, "WHERE") {
			t.Error("Expected WHERE clause in generated SQL")
		}

		// Test with published filter
		published := true
		filters = SearchFilters{Published: &published}
		query = searchService.BuildDynamicQuery(baseQuery, filters)
		sql, _, err = query.ToSql()
		if err != nil {
			t.Errorf("BuildDynamicQuery with published filter failed: %v", err)
		}
		if !strings.Contains(sql, "published") {
			t.Error("Expected published in generated SQL")
		}

		// Test with multiple filters
		filters = SearchFilters{
			Query:     "test",
			Published: &published,
			Limit:     10,
		}
		query = searchService.BuildDynamicQuery(baseQuery, filters)
		sql, _, err = query.ToSql()
		if err != nil {
			t.Errorf("BuildDynamicQuery with multiple filters failed: %v", err)
		}
		if !strings.Contains(sql, "LIMIT") {
			t.Error("Expected LIMIT clause in generated SQL")
		}
	})
}

// TestSquirrelQueryBuilder tests Squirrel query building functionality
func TestSquirrelQueryBuilder(t *testing.T) {
	// Test Squirrel query builder patterns
	t.Run("Basic Query Building", func(t *testing.T) {
		psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
		query := psql.Select("id", "name").From("users").Where(squirrel.Eq{"active": true})
		sql, args, err := query.ToSql()
		if err != nil {
			t.Errorf("Basic query building failed: %v", err)
		}
		if !strings.Contains(sql, "SELECT id, name FROM users WHERE active =") {
			t.Error("Expected proper SQL structure")
		}
		if len(args) != 1 || args[0] != true {
			t.Error("Expected correct arguments")
		}
	})

	t.Run("Complex Query Building", func(t *testing.T) {
		psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

		// Test JOINs
		query := psql.Select("u.id", "u.name", "p.title").
			From("users u").
			Join("posts p ON u.id = p.user_id").
			Where(squirrel.Eq{"p.published": true})

		sql, _, err := query.ToSql()
		if err != nil {
			t.Errorf("Complex query building failed: %v", err)
		}
		if !strings.Contains(sql, "JOIN") {
			t.Error("Expected JOIN clause")
		}

		// Test OR conditions
		query = psql.Select("*").From("posts").Where(squirrel.Or{
			squirrel.Eq{"published": true},
			squirrel.Eq{"user_id": 1},
		})

		sql, _, err = query.ToSql()
		if err != nil {
			t.Errorf("OR condition query failed: %v", err)
		}
		if !strings.Contains(sql, "OR") {
			t.Error("Expected OR condition")
		}
	})
}

// BenchmarkSquirrelVsManualSQL benchmarks Squirrel vs manual SQL building
func BenchmarkSquirrelVsManualSQL(b *testing.B) {
	// Compare performance of Squirrel vs manual string building
	b.Run("Squirrel", func(b *testing.B) {
		psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
		for i := 0; i < b.N; i++ {
			query := psql.Select("id", "name").From("users").Where(squirrel.Eq{"active": true})
			_, _, err := query.ToSql()
			if err != nil {
				b.Fatalf("Squirrel query failed: %v", err)
			}
		}
	})

	b.Run("Manual SQL", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sql := "SELECT id, name FROM users WHERE active = ?"
			args := []interface{}{true}
			_ = sql
			_ = args
		}
	})
}
