package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// Category represents a blog post category using GORM model conventions
// This model demonstrates GORM ORM patterns and relationships
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null;uniqueIndex"`
	Description string         `json:"description" gorm:"size:500"`
	Color       string         `json:"color" gorm:"size:7"` // Hex color code
	Active      bool           `json:"active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete support

	// GORM Associations (demonstrates ORM relationships)
	Posts []Post `json:"posts,omitempty" gorm:"many2many:post_categories;"`
}

// CreateCategoryRequest represents the payload for creating a category
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=500"`
	Color       string `json:"color" validate:"omitempty,hexcolor"`
}

// UpdateCategoryRequest represents the payload for updating a category
type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Color       *string `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Active      *bool   `json:"active,omitempty"`
}

// TableName specifies the table name for GORM (optional - GORM auto-infers)
func (Category) TableName() string {
	return "categories"
}

// BeforeCreate hook - GORM BeforeCreate hook
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	// Validate data before creation
	if c.Name == "" {
		return gorm.ErrInvalidData
	}

	// Set default values
	if c.Color == "" {
		c.Color = "#007bff"
	}

	// Set active to true by default
	c.Active = true

	// Log creation
	log.Printf("Creating category: %s", c.Name)

	return nil
}

// AfterCreate hook - GORM AfterCreate hook
func (c *Category) AfterCreate(tx *gorm.DB) error {
	// Log creation
	log.Printf("Category created successfully: %s (ID: %d)", c.Name, c.ID)

	// In a real application, you might:
	// - Send notifications
	// - Update cache
	// - Trigger webhooks

	return nil
}

// BeforeUpdate hook - GORM BeforeUpdate hook
func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	// Validate changes
	if c.Name == "" {
		return gorm.ErrInvalidData
	}

	// Log update
	log.Printf("Updating category: %s (ID: %d)", c.Name, c.ID)

	return nil
}

// Validate method for CreateCategoryRequest
func (req *CreateCategoryRequest) Validate() error {
	// Basic validation
	if req.Name == "" || len(req.Name) < 2 || len(req.Name) > 100 {
		return gorm.ErrInvalidData
	}

	// Description validation
	if len(req.Description) > 500 {
		return gorm.ErrInvalidData
	}

	// Color validation (basic hex check)
	if req.Color != "" && (len(req.Color) != 7 || req.Color[0] != '#') {
		return gorm.ErrInvalidData
	}

	return nil
}

// ToCategory method - Convert request to GORM model
func (req *CreateCategoryRequest) ToCategory() *Category {
	return &Category{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		Active:      true,
	}
}

// GORM scopes (reusable query logic)
func ActiveCategories(db *gorm.DB) *gorm.DB {
	return db.Where("active = ?", true)
}

func CategoriesWithPosts(db *gorm.DB) *gorm.DB {
	return db.Joins("Posts").Where("posts.id IS NOT NULL")
}

// Model validation methods
func (c *Category) IsActive() bool {
	return c.Active
}

func (c *Category) PostCount(db *gorm.DB) (int64, error) {
	count := db.Model(c).Association("Posts").Count()
	return count, nil
}
