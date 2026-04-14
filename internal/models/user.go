package models

import (
	roles "github.com/AsmrS4/certificates-plugin/internal/enums"
)

type (
	UserRole = roles.UserRole
)

type User struct {
	ID   int64    `json:"user_id"`
	Name string   `json:"username"`
	Role UserRole `json:"role"`
}
