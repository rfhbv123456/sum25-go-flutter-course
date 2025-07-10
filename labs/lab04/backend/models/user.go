package models

import (
	"database/sql"
	"errors"
	"regexp"
	"time"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest represents the payload for creating a user
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateUserRequest represents the payload for updating a user
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// Validate method for User
func (u *User) Validate() error {
	if len(u.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	if u.Email == "" || !isValidEmail(u.Email) {
		return errors.New("invalid email")
	}
	return nil
}

// Validate method for CreateUserRequest
func (req *CreateUserRequest) Validate() error {
	if len(req.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	if req.Email == "" || !isValidEmail(req.Email) {
		return errors.New("invalid email")
	}
	return nil
}

// ToUser method for CreateUserRequest
func (req *CreateUserRequest) ToUser() *User {
	now := time.Now()
	return &User{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// isValidEmail helper
func isValidEmail(email string) bool {
	// Very basic email regex
	re := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	return re.MatchString(email)
}

// TODO: Implement ScanRow method for User
func (u *User) ScanRow(row *sql.Row) error {
	// TODO: Scan database row into User struct
	// Handle the case where row might be nil
	return nil
}

// TODO: Implement ScanRows method for User slice
func ScanUsers(rows *sql.Rows) ([]User, error) {
	// TODO: Scan multiple database rows into User slice
	// Make sure to close rows and handle errors properly
	return nil, nil
}
