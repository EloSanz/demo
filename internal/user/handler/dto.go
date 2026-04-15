package handler

import (
	"github.com/elosanz/demo/internal/user"
)

// CreateUserRequest is the request body for creating a user.
type CreateUserRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
}

// UpdateUserRequest is the request body for updating a user.
type UpdateUserRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
}

// UserResponse is the API representation of a user.
type UserResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
}

// PageResponse is a paginated response wrapper for users.
type PageResponse struct {
	Content []UserResponse `json:"content"`
	Page    int            `json:"page"`
	Size    int            `json:"size"`
	Total   int64          `json:"total"`
}

func (r CreateUserRequest) toDomain() user.User {
	return user.User{Name: r.Name, Email: r.Email, Phone: r.Phone, Website: r.Website}
}

func (r UpdateUserRequest) toDomain() user.User {
	return user.User{Name: r.Name, Email: r.Email, Phone: r.Phone, Website: r.Website}
}

func mapToResponse(u user.User) UserResponse {
	return UserResponse{ID: u.ID, Name: u.Name, Email: u.Email, Phone: u.Phone, Website: u.Website}
}
