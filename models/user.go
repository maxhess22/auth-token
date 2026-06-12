package models

import "time"

type User struct {
	ID        int       `json:"id"`
	RoleId    string    `json:"role_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthInput struct {
	RoleId   string `json:"role_id"`
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
