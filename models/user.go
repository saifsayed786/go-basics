package models

import "time"

type UserRole struct {
	Id   int    `json:"id,string"`
	Role string `json:"role"`
}

type User struct {
	ID         int    `json:"id,string"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	UserRoleID int    `json:"role_id,string"`
	Password   string `json:"password"`

	// Meta
	CreatedAt time.Time `json:"created_at,timestamp"`
	UpdatedAt time.Time `json:"updated_at,timestamp"`
	IsActive  bool      `json:"is_active,boolean"`
}
