// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graph_model

import (
	"time"
)

type Account struct {
	ID        string    `json:"id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AccountAdd struct {
	Name     string `json:"name"`
	Gender   int    `json:"gender"`
	Birthday string `json:"birthday"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     int    `json:"type"`
}

type AccountChangePassword struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	NewPassword     string `json:"new_password"`
}

type AccountDelete struct {
	AccountID string `json:"account_id"`
}

type AccountPagination struct {
	Rows   []Account  `json:"rows"`
	Paging Pagination `json:"paging"`
}

type AccountSetPassword struct {
	AccountID   string `json:"account_id"`
	NewPassword string `json:"new_password"`
}

type AccountUpdate struct {
	AccountID *string    `json:"account_id,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Gender    *int       `json:"gender,omitempty"`
	Birthday  *time.Time `json:"birthday,omitempty"`
	Phone     *string    `json:"phone,omitempty"`
	Avatar    *string    `json:"avatar,omitempty"`
	Status    *int       `json:"status,omitempty"`
	Type      *int       `json:"type,omitempty"`
}

type DefaultResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	TotalPage   int `json:"total_page"`
	Total       int `json:"total"`
}