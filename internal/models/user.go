package models

import (
	enums "github.com/AsmrS4/certificates-plugin/internal/enums/role"
)

type (
	UserRole = enums.UserRole
)

type User struct {
	ID   int64    `json:"user_id"`
	Name string   `json:"username"`
	Role UserRole `json:"role"`
}
